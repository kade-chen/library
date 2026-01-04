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
func (t *Token) SetCookie(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		// Cookie 名称
		// 浏览器中通过 document.cookie 能看到 key，但拿不到 value（HttpOnly=true）
		Name: "CCK",
		// Cookie 的值
		// ⚠️ 你这里放的是 RefreshToken（名字和内容不一致，建议改名）
		Value: t.RefreshToken,
		// Cookie 的存活时间（秒）
		// 0  => Session Cookie（浏览器关闭即失效）
		// >0 => 明确秒数
		// <0 => 立即删除 Cookie
		MaxAge: 0,
		// Cookie 生效路径
		// "/" 表示整个站点所有路径都携带
		Path: "/",
		// Cookie 生效域名
		// ""  => 当前访问域名（如 api.wondercloud.com）
		// ".wondercloud.com" => 所有子域共享（大厂常用）
		Domain: "",
		// SameSite 策略（防 CSRF 核心参数）
		// DefaultMode ≈ Lax（现代浏览器）
		// Strict => 最严格
		// None   => 允许跨站（必须 Secure=true）
		SameSite: http.SameSiteLaxMode,
		// 是否只在 HTTPS 下发送
		// false => HTTP 也会带（⚠️ 生产环境非常不安全）
		// true  => 仅 HTTPS（大厂必开）
		Secure: false,
		// 是否禁止 JS 访问
		// true  => document.cookie 无法读取（防 XSS）
		// false => JS 可读（非常危险）
		HttpOnly: true,
	})
	return nil
}

// delete the cookie for web-ui
func (c *Token) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "CCK",           // 要清除的 Cookie 名称
		Value:   "",              // Cookie 的值为空
		Path:    "/",             // Cookie 的路径
		Expires: time.Unix(0, 0), // 过期时间设置为 Unix 纪元时间
		MaxAge:  -1,              // 立即失效
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
