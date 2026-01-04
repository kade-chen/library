package impl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/configs"
	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
	tools "github.com/kade-chen/google-billing-console/tools/rand"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
)

func (t *service) GeneratJwtAccessToken(platform int32, subject string, issueAt, ExpiredAt int64, domains []string) (string, error) {
	id, err := tools.NewJwtId()
	if err != nil {
		return "", err
	}
	claims := authModel.TokenAuthMiddleware{
		Platform: int32(platform),
		Scope:    []string{"platform.admin"},
		Domains:  domains,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "wondercloud.com",
			Subject:   subject,
			Audience:  jwt.ClaimStrings{"dev.billing.wondercloud.com"},
			IssuedAt:  jwt.NewNumericDate(time.Unix(issueAt, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(issueAt, 0).Add(time.Duration(ExpiredAt) * time.Second)),
			ID:        id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	cc, err := token.SignedString(ioc.Config().Get(configs.AppName).(*config.Service).JwtPrivateKey)
	if err != nil {
		return "", exception.NewInternalServerError("generate token failed")
	}
	return cc, nil
}
