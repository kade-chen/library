package project

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "project"
)

type Service interface {
	// QueryByProject(ctx context.Context, query string) error
	QueryByDateProjectAll(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByDateProjectServicesCustomSku(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByDateProjectCustomServicesAllSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByDateProjectCustomServicesCustomSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error)
	// QueryBySku(ctx context.Context, query string) error
	// QueryByDateSku(ctx context.Context, query string) error
	// QueryByService(ctx context.Context, query string) error
	// QueryByDateService(ctx context.Context, query string) error
}
