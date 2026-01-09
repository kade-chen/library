package organization

import (
	"github.com/go-playground/validator/v10"
	"github.com/kade-chen/library/exception"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New() // 创建一个验证器
)

// verity that the request is valid
func (req *CreateOrganizationRequest) validate_strust() error {
	return validate.Struct(req)
}

// verity that the request is valid
func (req *DescribeOrganizationRequest) Validate() error {
	switch req.DescribeBy {
	case DESCRIBE_BY_ID:
		if req.Id == "" {
			return exception.NewBadRequest("Organization id required")
		}
	case DESCRIBE_BY_NAME:
		if req.Name == "" {
			return exception.NewBadRequest("Organization name required")
		}
	}

	return validate.Struct(req)
}
