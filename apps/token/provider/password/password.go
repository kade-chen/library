package password

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

var (
	AUTH_FALIED = exception.NewUnauthorized("user or password not empty")
)

// var _ provider.Issuer = &password{}

var _ provider.Issuer = (*password)(nil)

func init() {
	provider.Registe(&password{})
}

type password struct {
	bq_client *bigquery.Client
	user      user.Service
	log       *zerolog.Logger
	domain    domain.Service
}

func (p *password) Init() error {
	p.bq_client = ioc.Config().Get(configs.AppName).(*impl.Service).BQ
	p.user = ioc.Controller().Get(user.AppName).(user.Service)
	p.domain = ioc.Controller().Get(domain.AppName).(domain.Service)
	p.log = log.Sub("issuer_password_token")
	return nil
}

func (*password) GrantType() token.GRANT_TYPE {
	return token.GRANT_TYPE_PASSWORD // 密码令牌
}

func (p *password) IssueToken(ctx context.Context, req *token.IssueTokenRequest) (*token.Token, error) {
	p.log.Info().Msg("The token issuance method is follows: PASSWORD")
	u, err := p.validate(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	// 3. 颁发Token
	tk := token.NewToken(req)
	tk.Domain = u.Spec.Domain
	tk.SharedUser = u.Spec.Shared
	tk.Username = u.Spec.Username
	tk.UserType = u.Spec.Type
	tk.UserId = u.Id

	return tk, err
}

func (p *password) validate(ctx context.Context, username, password string) (*user.User, error) {
	if username == "" || password == "" {
		p.log.Error().Msg("username or password is empty")
		return nil, AUTH_FALIED
	}

	//1.query user whether in the document
	p.log.Info().Msgf("query user by username: %s", username)
	p.log.Info().Msgf("query user by password: %s", password)
	u, err := p.user.DescribeUser(ctx, &user.DescribeUserRequest{
		DescribeBy: user.DESCRIBE_BY_USER_ID,
		Id:         username,
	})
	if err != nil {
		return nil, err
	}
	//2.verify password
	if err := u.Password.CheckPassword(password); err != nil {
		return nil, err
	}

	// 2.check whether the password has expried
	var expiredRemain, expiredDays uint
	switch u.Spec.Type {
	case user.TYPE_SUB:
		p.log.Info().Msg("sub account")
		// 子账号过期策略
		d, err := p.domain.DescribeDomain(ctx, domain.NewDescribeDomainRequestByName(u.Spec.Domain))
		if err != nil {
			return nil, err
		}
		ps := d.Spec.PasswordConfig
		//BeforeExpiredRemindDays=10  password_expired_days=90
		expiredRemain, expiredDays = uint(ps.BeforeExpiredRemindDays), uint(ps.PasswordExpiredDays)
	case user.TYPE_PRIMARY:
		p.log.Info().Msg("primary user")
		return u, nil
	case user.TYPE_SUPPER:
		p.log.Info().Msg("super user")
		return u, nil
	default:
		// 主账号和管理员密码过期策略
		expiredRemain, expiredDays = uint(u.Password.ExpiredRemind), uint(u.Password.ExpiredDays)
	}
	err = u.Password.CheckPasswordExpired(ctx, expiredRemain, expiredDays, p.bq_client, u.Id)
	if err != nil {
		p.log.Error().Msgf("check password expired error: %s", err)
		return nil, err
	}
	return u, err
}
