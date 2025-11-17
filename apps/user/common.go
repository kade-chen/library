package user

import (
	"github.com/go-playground/validator"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)
