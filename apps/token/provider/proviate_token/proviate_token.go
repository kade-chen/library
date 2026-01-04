package proviatetoken

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

// 接口约束
var _ provider.Issuer = (*proviatetoken)(nil)

type proviatetoken struct {
	log   *zerolog.Logger
	token token.Service
}

func init() {
	provider.Registe(&proviatetoken{})

}

func (i *proviatetoken) GrantType() token.GRANT_TYPE {
	return token.GRANT_TYPE_PRIVATE_TOKEN // 私有令牌
}

func (i *proviatetoken) Init() error {
	i.token = ioc.Controller().Get(token.AppName).(token.Service)
	i.log = log.Sub("issuer.private_token")
	return nil
}

func (i *proviatetoken) IssueToken(ctx context.Context, req *token.IssueTokenRequest) (*token.Token, error) {
	// i.Init()
	//0.verity accessToken
	i.log.Info().Msg("verity access_token whether empty")
	tk, err := i.verityAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	i.log.Info().Msg("verity access_token is finsh")
	//1.颁发token
	i.log.Info().Msg("颁发GRANT_TYPE_PRIVATE_TOKEN program doing")
	token := token.NewToken(req)
	token.AccessToken = tk.AccessToken
	i.log.Info().Msg("颁发GRANT_TYPE_PRIVATE_TOKEN program success")
	return token, nil
}

func (i *proviatetoken) verityAccessToken(ctx context.Context, accessToken string) (*token.Token, error) {
	//0.verity accessToken whether empty
	i.log.Info().Msg("start veritication accessToken")
	if accessToken == "" {
		return nil, exception.NewNotFound("accessToken is required")
	}
	// 1.判断凭证合法性
	i.log.Info().Msg("Legitimacy of judgment for accessToken")
	token, err := i.token.ValicateToken(ctx, token.NewValidateTokenRequest())
	if err != nil {
		return nil, err
	}

	return token, nil
}
