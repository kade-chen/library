package rand

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/kade-chen/library/exception"
)

func NewJwtId() (string, error) {
	b := make([]byte, 16) // 128 bit
	if _, err := rand.Read(b); err != nil {
		return "", exception.NewInternalServerError("rand: failed to read random bytes: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
