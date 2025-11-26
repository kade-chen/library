package usagedate

import (
	"context"

	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
)

const (
	AppName = "usagedate"
)

type ProjectService interface {
	//辅助功能
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataServiceSkusRequest) (model.ByDateProjectAllServicesSkusList, error)
	//by date project
	QueryByDateProject(context.Context, *model.ProjectDataRequest) ([]model.ProjectDateCost, error)

	//by project
	QueryByProject(context.Context, *model.ProjectRequest) ([]model.ProjectCost, error)
}

type Service interface {
	//辅助功能
	QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataServiceSkusRequest) ([]model.ServicesList, error)
	//by date Service
	QueryByDateService(context.Context, *model.ServiceDataRequest) ([]model.ServiceDateCost, error)

	// //by Service
	QueryByService(context.Context, *model.ServiceRequest) ([]model.ServiceCost, error)
}

type SkuService interface {
	QueryByDateProjectSKUsAll(ctx context.Context, config *model.ProjectDataServiceSkusRequest) ([]model.SkusList, error)

	//by date sku
	QueryByDateSku(context.Context, *model.SkuDataRequest) ([]model.SkuDateCost, error)

	// //by Service
	QueryBySku(context.Context, *model.SkuRequest) ([]model.SkuCost, error)
}
