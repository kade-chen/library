package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type TokenAuthMiddleware struct {
	// user user.Service
	TrancesID     string   `json:"trances_id"` // 请求ID
	Platform      int32    `json:"platform"`   // web / sdk / admin
	Scope         []string `json:"scope"`
	JwtToken      string   `json:"jwt_token"`
	Organizations []string `json:"organizations"`
	// role user.Role
	jwt.RegisteredClaims
}
