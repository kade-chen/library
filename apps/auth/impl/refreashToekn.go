package impl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/configs"
	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
	tools "github.com/kade-chen/google-billing-console/tools/rand"
	"github.com/kade-chen/library/ioc"
)

func (t *service) JwtRefreshAccessToken(platform int32, subject string, expiredAt int64, organizations []string) (accesstoken string, err error) {
	id, err := tools.NewJwtId()
	if err != nil {
		return "", err
	}

	claims := authModel.TokenAuthMiddleware{
		Platform:      int32(platform),
		Scope:         []string{"platform.admin"},
		Organizations: organizations,
		// Type: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "wondercloud.com",
			Subject:   subject,
			Audience:  jwt.ClaimStrings{"dev.billing.wondercloud.com"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredAt) * time.Second)),
			ID:        id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(
		ioc.Config().Get(configs.AppName).(*config.Service).JwtPrivateKey,
	)
}
