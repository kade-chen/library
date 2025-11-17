package token

import context "context"

const (
	AppName = "token"
)

type Service interface {
	//issue an token for domain
	IssueToken(context.Context, *IssueTokenRequest) (*Token, error)
	// remove Token
	RevolkToken(context.Context, *RevolkTokenRequest) (*Token, error)
	// RPC
	RPCServer
}
