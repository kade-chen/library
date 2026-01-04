package middlewares

// import (
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/kade-chen/google-billing-console/apps/configs"
// 	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
// 	tools "github.com/kade-chen/google-billing-console/tools/rand"
// 	"github.com/kade-chen/library/ioc"
// )

// func GenerateRefreshToken(platform int32, subject string, expiredAt int64) (accesstoken string, err error) {
// 	id, err := tools.NewJwtId()
// 	if err != nil {
// 		return "", err
// 	}

// 	claims := TokenAuthMiddleware{
// 		Platform: int32(platform),
// 		Scope:    []string{"platform.admin"},
// 		// Type: TokenTypeRefresh,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Issuer:    "wondercloud.com",
// 			Subject:   subject,
// 			Audience:  jwt.ClaimStrings{"dev.billing.wondercloud.com"},
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredAt) * time.Second)),
// 			ID:        id,
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	return token.SignedString(
// 		ioc.Config().Get(configs.AppName).(*config.Service).JwtPrivateKey,
// 	)
// }
