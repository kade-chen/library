package middlewares

// import (
// 	"errors"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/emicklei/go-restful/v3"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/kade-chen/google-billing-console/apps/configs"
// 	config "github.com/kade-chen/google-billing-console/apps/configs/impl"
// 	"github.com/kade-chen/google-billing-console/apps/token"
// 	"github.com/kade-chen/library/exception"
// 	"github.com/kade-chen/library/http/response"
// 	"github.com/kade-chen/library/ioc"
// 	logs "github.com/kade-chen/library/ioc/config/log"
// 	"github.com/rs/zerolog"
// )

// // 用于鉴权的中间件
// // 用于Token鉴权的中间件
// type TokenAuthMiddleware struct {
// 	tk token.Service
// 	// user user.Service
// 	Platform int32    `json:"platform"` // web / sdk / admin
// 	Scope    []string `json:"scope"`
// 	JwtToken string   `json:"jwt_token"`
// 	log      *zerolog.Logger
// 	// ioc.ObjectImpl
// 	// role user.Role
// 	jwt.RegisteredClaims
// }

// func NewTokenAuthMiddleware() *TokenAuthMiddleware {
// 	return &TokenAuthMiddleware{
// 		// tk: ioc.Default().Get(tokenimpl.AppName).(*tokenimpl.TokenServiceImpl),
// 		tk: ioc.Controller().Get(token.AppName).(token.Service),
// 		// user: ioc.Controller().Get(user.AppName).(user.Service),
// 		log: logs.Sub("TokenAuthMiddleware"),
// 	}
// }

// func (t *TokenAuthMiddleware) Auth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 	auth := req.HeaderParameter("Authorization")
// 	if !strings.HasPrefix(auth, "Bearer ") {
// 		response.Failed(resp, exception.NewUnauthorized("missing bearer token"))
// 		// resp.WriteErrorString(http.StatusUnauthorized, "missing bearer token")
// 		return
// 	}

// 	tokenStr := strings.TrimPrefix(auth, "Bearer ")
// 	t.JwtToken = tokenStr
// 	claims, err := t.parseToken()
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, jwt.ErrTokenExpired):
// 			// access token 过期 → refresh
// 			Redirect_Url(resp, err, t.log)

// 		case errors.Is(err, jwt.ErrSignatureInvalid):
// 			// 伪造 token
// 			response.Failed(resp, exception.NewUnauthorized("invalid signature"))

// 		default:
// 			response.Failed(resp, exception.NewUnauthorized("invalid token"))
// 		}
// 		return
// 	}

// 	// ❗❗❗ 平台隔离（你问的这一行）
// 	if claims.Platform != 0 {
// 		response.Failed(resp, exception.NewForbidden("platform not allowed"))
// 		return
// 	}

// 	// 3.wether the token is exist
// 	_, err = t.tk.ValicateToken(req.Request.Context(), &token.ValicateTokenRequest{
// 		AccessToken: tokenStr,
// 	})
// 	if err != nil {
// 		// 处理 token 验证错误
// 		t.log.Error().Msgf("The %s is expired or not found", tokenStr)
// 		Redirect_Url(resp, err, t.log)
// 		return
// 	}

// 	// 存入 request attribute
// 	// fmt.Println(format.ToJSON(claims))
// 	req.SetAttribute("claims", claims)
// 	// v := req.Attribute("claims").(string)
// 	chain.ProcessFilter(req, resp)
// }

// func (t *TokenAuthMiddleware) parseToken() (*TokenAuthMiddleware, error) {
// 	// 使用 jwt 库解析 tokenStr，目标是解析到 Claims 结构体中
// 	jwttoken, err := jwt.ParseWithClaims(t.JwtToken, &TokenAuthMiddleware{}, func(token *jwt.Token) (interface{}, error) {
// 		// 安全校验：防止 alg 被篡改
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, jwt.ErrSignatureInvalid
// 		}

// 		// 提供签名密钥，用于验证 JWT 的签名是否合法
// 		return ioc.Config().Get(configs.AppName).(*config.Service).JwtPublicKey, nil
// 	},
// 		jwt.WithAudience("dev.billing.wondercloud.com"),
// 		jwt.WithIssuer("wondercloud.com"),
// 	)

// 	if jwttoken != nil {
// 		if claims, ok := jwttoken.Claims.(*TokenAuthMiddleware); ok && claims.ExpiresAt != nil {
// 			t.log.Info().Msgf("ID: %v, Issuer: %v, Audience: %v, Platform: %v, Scope: %v, Subject:%v, Token: %v, Expiration Status: %v, Expiration Time: %v", claims.ID, claims.Issuer, claims.Audience, token.PLATFORM(int32(claims.Platform)), claims.Scope, claims.Subject, t.JwtToken, claims.ExpiresAt.Time.Before(time.Now()), claims.ExpiresAt.Time)
// 		}
// 	}

// 	// 如果解析失败或 token 无效，则返回错误
// 	if err != nil {
// 		return nil, err // ❗原样返回
// 	}

// 	// 将解析出的 Claims 转换为自定义的 *Claims 类型并返回
// 	return jwttoken.Claims.(*TokenAuthMiddleware), nil
// }

// func Redirect_Url(resp *restful.Response, err error, t *zerolog.Logger) {
// 	redirectURL := "http://localhost:5173/login?error=" + err.Error()
// 	resp.AddHeader("Location", redirectURL)
// 	resp.WriteHeader(http.StatusFound) // 302
// 	t.Warn().Msg("302 to login page")
// }
