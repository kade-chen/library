package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/tools/format"
)

func (t *Token) Json() string {
	return format.ToJSON(t)
}

// 基于令牌创建HTTP Cookie 用于Web登陆场景
func (c *Token) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "ACCESS_TOKEN_COOKIE_KEY",
		Value:    c.AccessToken,
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	})
}

// delete the cookie for web-ui
func (c *Token) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "ACCESS_TOKEN_COOKIE_KEY", // 要清除的 Cookie 名称
		Value:   "",                        // Cookie 的值为空
		Path:    "/",                       // Cookie 的路径
		Expires: time.Unix(0, 0),           // 过期时间设置为 Unix 纪元时间
		MaxAge:  -1,                        // 立即失效
	})
}

// // Equal type compare
// func (t PLATFORM) Equal(target PLATFORM) bool {
// 	return t == target
// }

// check issue_token expired
func (t *Token) CheckIssue_Token_expried(issue_at, access_expired_at int64, AccessToken string) (string, error) {
	access_expired := time.Unix(issue_at, 0).Add(time.Duration(access_expired_at) * time.Second)

	expired_time := time.Since(access_expired).Seconds()
	// fmt.Println("expired_time:", expired_time, access_expired)
	if expired_time > 0 {
		return "", exception.NewAccessTokenExpired("access token(%s) expried, expired_time:%s ,expried_at: %v Second", AccessToken, access_expired, expired_time)
	}

	return fmt.Sprintf("access token:%s, expired_time:%s ,expried_at: %v Second", AccessToken, access_expired, -expired_time), nil
}
