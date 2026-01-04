package token

import (
	context "context"
)

const (
	AppName = "token"
)

type Service interface {
	//issue an token for domain
	IssueToken(context.Context, *IssueTokenRequest) (*Token, error)
	//update an token for domain
	UpdateToken(context.Context, *UpdateTokenRequest) (string, error)
	// remove Token
	RevolkToken(context.Context, *RevolkTokenRequest) (*Token, int64, error)
	// RPC
	RPCServer
}
