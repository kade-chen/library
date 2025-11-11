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
	QueryByProjectAllQueryByDateProjectAll(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectServicesCustomSku(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectCustomServicesAllSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectCustomServicesCustomSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error)
	// QueryBySku(ctx context.Context, query string) error

	QueryByProjectAll(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	// QueryByDateSku(ctx context.Context, query string) error
	// QueryByService(ctx context.Context, query string) error
	// QueryByDateService(ctx context.Context, query string) error
}
