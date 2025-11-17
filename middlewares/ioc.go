package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	logs "github.com/kade-chen/library/ioc/config/log"

	"github.com/emicklei/go-restful/v3"
	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/rs/zerolog"
)

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type TokenAuth struct {
	tk   token.Service
	user user.Service
	log  *zerolog.Logger
	// role user.Role
}

func NewTokenAuther() *TokenAuth {
	return &TokenAuth{
		// tk: ioc.Default().Get(tokenimpl.AppName).(*tokenimpl.TokenServiceImpl),
		tk:   ioc.Controller().Get(token.AppName).(token.Service),
		user: ioc.Controller().Get(user.AppName).(user.Service),
		log:  logs.Sub("TokenAuth"),
	}
}

func (t *TokenAuth) Auth_Login(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		t.log.Error().Msgf("Internal Server Error; error: %v", err)
	// 		response.Failed(resp, exception.NewInternalServerError("Internal Server Error; error: %v", err))
	// 		return
	// 	}
	// }()

	//	1.get cookis
	http_cookis := req.Request.Cookies()
	if len(http_cookis) == 0 {
		t.log.Info().Msg("No cookies found")
		Redirect_Url(resp, exception.NewAccessTokenIllegal("No cookies found"), t.log)
		return
	}
	var cookie_name, cookie_value string
	for _, cookie := range http_cookis {
		cookie_name = cookie.Name
		cookie_value = cookie.Value
	}
	//	2.wether the cookie is exist
	if cookie_name != "ACCESS_TOKEN_COOKIE_KEY" {
		t.log.Info().Msg("COOKIES_KEY NOT FOUND")
		resp.AddHeader("Location", "http://localhost:5173/login")
		resp.WriteHeader(http.StatusFound) // 302
		return
	}
	// 3.wether the token is exist
	tk, err := t.tk.ValicateToken(req.Request.Context(), &token.ValicateTokenRequest{
		AccessToken: cookie_value,
	})
	if err != nil {
		// 处理 token 验证错误
		t.log.Error().Msgf("The %s is expired or not found", cookie_value)
		Redirect_Url(resp, err, t.log)
		return
	}
	// 4.query wether the user exist
	user, err := t.user.DescribeUser(req.Request.Context(), &user.DescribeUserRequest{
		DescribeBy: user.DESCRIBE_BY_USER_NAME,
		Username:   tk.Username,
	})
	if err != nil {
		t.log.Error().Msgf("The %s not found", tk.Username)
		// 处理 user 验证错误
		Redirect_Url(resp, err, t.log)
		return
	}
	// 5.wether the user's domain is exist
	if tk.Domain != user.Spec.Domain {
		t.log.Error().Msgf("The %s not found", tk.Domain)
		// 处理 token 验证错误
		Redirect_Url(resp, exception.NewAccessTokenIllegal("the domain not find for user"), t.log)
		return
	}
	// 6.verity whether the issue_token is expired
	expried_access_token_time, expried_time, err := NewCheckIssue_Token(tk)
	if err != nil {
		t.log.Error().Msg("The issue_token is expired")
		// 处理 token 验证错误
		Redirect_Url(resp, err, t.log)
		return
	}
	t.log.Info().Msgf("the token is valid, expried_access_token_time: %v , expried_time: %v Second", expried_access_token_time, -expried_time)

	// 7.set the token to the context
	ctx := context.WithValue(req.Request.Context(), "user", user)
	// req.Request = req.Request.WithContext(ctx)
	// ctx = context.WithValue(req.Request.Context(), "useruname", user.Spec.Username)
	req.Request = req.Request.WithContext(ctx)
	chain.ProcessFilter(req, resp)
}

func Redirect_Url(resp *restful.Response, err error, t *zerolog.Logger) {
	redirectURL := "http://localhost:5173/login?error=" + err.Error()
	resp.AddHeader("Location", redirectURL)
	resp.WriteHeader(http.StatusFound) // 302
	t.Error().Msg("302 to login page")
}

func NewCheckIssue_Token(tk *token.Token) (time.Time, float64, error) {
	expried_access_token_time := time.Unix(tk.IssueAt, 0).Add(time.Duration(tk.AccessExpiredAt) * time.Second)
	expried_time := time.Since(expried_access_token_time).Seconds()
	if expried_time >= 0 {
		return expried_access_token_time, expried_time, exception.NewAccessTokenExpired("token %s 过期了 %f秒", tk.UserId, expried_time)
	}
	return expried_access_token_time, expried_time, nil
}
