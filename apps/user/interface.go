package user

import (
	"context"
)

const (
	AppName = "users"
)

type Service interface {
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*UserSet, error)
	RPCServer
}
