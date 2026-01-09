package organization

import "context"

const (
	AppName = "organizations"
)

type Service interface {
	// 创建域
	CreateOrganization(context.Context, *CreateOrganizationRequest) (*Organization, error)
	// 更新域
	// UpdateOrganization(context.Context, *UpdateOrganizationRequest) (*Organization, error)
	// RPC
	RPCServer
}
