package usagedate

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "usagedate"
)

type ProjectService interface {
	//辅助功能
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error)
	//by date project
	QueryByDateProject(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)

	//by project
	QueryByProject(context.Context, *model.ProjectConfig) ([]model.ProjectCost, error)
}

type Service interface {
	//辅助功能
	QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.ServicesList, error)
	//by date Service
	QueryByDateService(context.Context, *model.ServiceDataConfig) ([]model.ServiceDateCost, error)

	// //by Service
	QueryByService(context.Context, *model.ServiceConfig) ([]model.ServiceCost, error)
}

type SkuService interface {
	QueryByDateProjectSKUsAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.SkusList, error)

	//by date sku
	QueryByDateSku(context.Context, *model.SkuDataConfig) ([]model.SkuDateCost, error)

	// //by Service
	QueryBySku(context.Context, *model.SkuConfig) ([]model.SkuCost, error)
}
