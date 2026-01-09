package impl

import (
	"errors"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v5"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/configs"
	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/tools/trances"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
	"github.com/kade-chen/library/ioc"
)

func (t *service) Auth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	t.auth.TrancesID = trances.NewTraceID()
	auth := req.HeaderParameter("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		response.Failed(resp, exception.NewUnauthorized("missing bearer token"))
		// resp.WriteErrorString(http.StatusUnauthorized, "missing bearer token")
		return
	}

	t.auth.JwtToken = strings.TrimPrefix(auth, "Bearer ")
	claims, err := t.parseToken()
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			// access token 过期 → refresh
			// Redirect_Url(resp, err, t.log)
			response.Failed(resp, exception.NewUnauthorized("token expired"))
		case errors.Is(err, jwt.ErrSignatureInvalid):
			// 伪造 token
			response.Failed(resp, exception.NewUnauthorized("invalid signature"))

		default:
			response.Failed(resp, exception.NewUnauthorized("invalid token"))
		}
		return
	}

	// ❗❗❗ 平台隔离（你问的这一行）
	if claims.Platform != 0 {
		response.Failed(resp, exception.NewForbidden("platform not allowed"))
		return
	}

	// 3.wether the token is exist
	_, err = t.tk.ValicateToken(req.Request.Context(), &token.ValicateTokenRequest{
		AccessToken: t.auth.JwtToken,
	})
	if err != nil {
		// 处理 token 验证错误
		t.log.Error().Msgf("The %s is expired or not found, ERROR: %v", t.auth.JwtToken, err)
		response.Failed(resp, exception.NewUnauthorized("The %v is expired or not found", t.auth.JwtToken))
		// Redirect_Url(resp, err, t.log)
		return
	}

	// 存入 request attribute
	// fmt.Println(format.ToJSON(claims))
	req.SetAttribute("claims", claims)
	// v := req.Attribute("claims").(string)
	chain.ProcessFilter(req, resp)
}

func (t *service) parseToken() (*authModel.TokenAuthMiddleware, error) {
	// 使用 jwt 库解析 tokenStr，目标是解析到 Claims 结构体中
	jwttoken, err := jwt.ParseWithClaims(t.auth.JwtToken, &authModel.TokenAuthMiddleware{}, func(token *jwt.Token) (interface{}, error) {
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
	// Claims := jwttoken.Claims.(*authModel.TokenAuthMiddleware)
	if jwttoken != nil {
		if claims, ok := jwttoken.Claims.(*authModel.TokenAuthMiddleware); ok && claims.ExpiresAt != nil {
			claims.JwtToken = t.auth.JwtToken
			claims.TrancesID = t.auth.TrancesID
			t.log.Info().Msgf("ID: %v, TraceID: %v, Organizations: %v, Issuer: %v, Audience: %v, Platform: %v, Scope: %v, Subject:%v, Token: %v, Expiration Status: %v, Expiration Time: %v",
				claims.ID, claims.TrancesID, claims.Organizations, claims.Issuer, claims.Audience, token.PLATFORM(int32(claims.Platform)), claims.Scope, claims.Subject, claims.JwtToken, claims.ExpiresAt.Time.Before(time.Now()), claims.ExpiresAt.Time)
		}
	}

	// 如果解析失败或 token 无效，则返回错误
	if err != nil {
		return nil, err // ❗原样返回
	}

	// 将解析出的 Claims 转换为自定义的 *Claims 类型并返回
	return jwttoken.Claims.(*authModel.TokenAuthMiddleware), nil
}

// func Redirect_Url(resp *restful.Response, err error, t *zerolog.Logger) {
// 	redirectURL := "http://localhost:5173/login?error=" + err.Error()
// 	resp.AddHeader("Location", redirectURL)
// 	resp.WriteHeader(http.StatusFound) // 302
// 	t.Warn().Msg("302 to login page")
// }
