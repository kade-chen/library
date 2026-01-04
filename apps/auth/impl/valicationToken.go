package impl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/configs"
	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/ioc"
)

func (t *service) ValicateToken(jwtToken string) (*authModel.TokenAuthMiddleware, error) {
	// 使用 jwt 库解析 tokenStr，目标是解析到 Claims 结构体中
	jwttoken, err := jwt.ParseWithClaims(jwtToken, &authModel.TokenAuthMiddleware{}, func(token *jwt.Token) (interface{}, error) {
		// 安全校验：防止 alg 被篡改
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// 提供签名密钥，用于验证 JWT 的签名是否合法
		return ioc.Config().Get(configs.AppName).(*config.Service).JwtPublicKey, nil
	},
		jwt.WithAudience("dev.billing.wondercloud.com"),
		jwt.WithIssuer("wondercloud.com"),
	)

	if jwttoken != nil {
		// fmt.Println(format.ToJSON(jwttoken.Claims))
		if claims, ok := jwttoken.Claims.(*authModel.TokenAuthMiddleware); ok && claims.ExpiresAt != nil {
			claims.JwtToken = jwtToken
			t.log.Info().Msgf("ID: %v, Issuer: %v, Audience: %v, Platform: %v, Scope: %v, Subject:%v, Token: %v, Expiration Status: %v, Expiration Time: %v", claims.ID, claims.Issuer, claims.Audience, token.PLATFORM(int32(claims.Platform)), claims.Scope, claims.Subject, jwtToken, claims.ExpiresAt.Time.Before(time.Now()), claims.ExpiresAt.Time)
		}
	}

	// 如果解析失败或 token 无效，则返回错误
	if err != nil {
		return jwttoken.Claims.(*authModel.TokenAuthMiddleware), err // ❗原样返回
	}

	// 将解析出的 Claims 转换为自定义的 *Claims 类型并返回
	return jwttoken.Claims.(*authModel.TokenAuthMiddleware), nil
}
