package token

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

// Validate 校验参数
func (m *ValicateTokenRequest) VerifyStruct() error {
	// 1.verify sturst
	if err := validate.Struct(m); err != nil {
		return err
	}

	if m.AccessToken == "" {
		return errors.New("access_token required one")
	}

	return nil
}
