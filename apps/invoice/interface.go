package invoice

import (
	"context"

	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
)

const (
	AppName = "invoice"
)

type ProjectService interface {
	//辅助功能
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataServiceSkuRequest) (model.ByDateProjectAllServicesSkusList, error)
	//by date project
	QueryByDateProject(context.Context, *model.ProjectDataRequest) ([]model.ProjectDateCost, error)

	//by project
	QueryByProject(context.Context, *model.ProjectRequest) ([]model.ProjectCost, error)
}

type Service interface {
	//辅助功能
	QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataServiceSkuRequest) ([]model.ServicesList, error)
	//by date Service
	QueryByDateService(context.Context, *model.ServiceDataRequest) ([]model.ServiceDateCost, error)

	// //by Service
	QueryByService(context.Context, *model.ServiceRequest) ([]model.ServiceCost, error)
}

type SkuService interface {
	QueryByDateProjectSKUsAll(ctx context.Context, config *model.ProjectDataServiceSkuRequest) ([]model.SkusList, error)

	//by date sku
	QueryByDateSku(context.Context, *model.SkuDataRequest) ([]model.SkuDateCost, error)

	// //by Service
	QueryBySku(context.Context, *model.SkuRequest) ([]model.SkuCost, error)

	QueryByDateSkuHeru(context.Context, *model.SkuDataRequest) ([]model.AlibabaHehuSkuDateCost, error)
}
