package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New() // 创建一个验证器
)

// verity that the request is valid
func (req *CreateDomainRequest) validate_strust() error {
	return validate.Struct(req)
}

// verity that the request is valid
func (req *DescribeDomainRequest) Validate() error {
	switch req.DescribeBy {
	case DESCRIBE_BY_ID:
		if req.Id == "" {
			return fmt.Errorf("domain id required")
		}
	case DESCRIBE_BY_NAME:
		if req.Name == "" {
			return fmt.Errorf("domain name required")
		}
	}

	return validate.Struct(req)
}
