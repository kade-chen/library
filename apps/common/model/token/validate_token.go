package token

type ValidateToken struct {
	JwtToken string `json:"jwt_token"`

	EnableExpire bool `json:"enable_expire"`

	ExpiresAt string `json:"expires_at"`
}
