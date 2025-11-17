package token

import (
	"math/rand"
	"strings"
	"time"
)

func NewToken(req *IssueTokenRequest) *Token {
	tk := &Token{
		AccessToken:      NewAccessTokenAt(24),
		AccessExpiredAt:  req.ExpiredAt,
		RefreshToken:     NewAccessTokenAt(32),
		RefreshExpiredAt: req.ExpiredAt * 4,
		IssueAt:          time.Now().Unix(),
		// AccessToken:     req.AccessToken,
		Description: req.Description,
		GrantType:   req.GrantType,
		Type:        req.Type,
		Status:      &Status{
			// IsBlock: false,
		},
		Location: &Location{

			IpLocation: &IPLocation{},
			UserAgent:  &UserAgent{},
			// ...其他配置
		},
		Meta: map[string]string{},
	}
	switch tk.GrantType {
	case GRANT_TYPE_PRIVATE_TOKEN:
		tk.Platform = PLATFORM_API
	default:
		tk.Platform = PLATFORM_WEB
	}
	return tk
}

// -----------------
// NewIssueTokenRequest 默认请求
func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{
		ExpiredAt: DEFAULT_ACCESS_TOKEN_EXPIRE_SECOND,
	}
}

// privateToken
func NewPrivateIssueTokenRequest(accessToken, describe string) *IssueTokenRequest {
	req := NewIssueTokenRequest()
	req.GrantType = GRANT_TYPE_PRIVATE_TOKEN
	req.AccessToken = accessToken
	req.Description = describe
	return req
}

func NewAccessTokenAt(lens int) string {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	t := make([]string, lens)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(lens) + rand.Int63n(1000)))
	for i := 0; i < lens; i++ {
		rint := r.Intn(len(charlist))
		p := charlist[rint : rint+1]
		t = append(t, p)
	}
	temporarily := strings.Join(t, "")
	return temporarily
}

func NewValidateTokenRequest(accessToken string) *ValicateTokenRequest {
	return &ValicateTokenRequest{
		AccessToken: accessToken,
	}
}

// --------------------

func NewStatus() *Status {
	return &Status{
		IsBlock: false,
	}
}

// NewRevolkTokenRequest 撤销Token请求
func NewRevolkTokenRequest(accessToken, refreshToken string) *RevolkTokenRequest {
	return &RevolkTokenRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// 优先从认证头中获取, 如果头没有从Cookie中获取
// func GetAccessTokenFromHTTP(r *http.Request) string {
// 	auth := r.Header.Get(ACCESS_TOKEN_HEADER_KEY)
// 	info := strings.Split(auth, " ")
// 	if len(info) > 1 {
// 		return info[1]
// 	}

// 	ck, err := r.Cookie(ACCESS_TOKEN_COOKIE_KEY)
// 	if err != nil {
// 		log.L().Warn().Msgf("get tk from cookie: %s error, %s", ACCESS_TOKEN_COOKIE_KEY, err)
// 		return r.URL.Query().Get(ACCESS_TOKEN_QUERY_KEY)
// 	}

// 	return ck.Value
// }
