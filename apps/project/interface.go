package project

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "project"
)

type Service interface {
	//by date project
	QueryByDateProjectAll(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectServicesCustomSku(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectCustomServicesAllSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectCustomServicesCustomSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error)

	//by project
	QueryByProjectAll(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByProjectServicesCustomSku(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByProjectCustomServicesAllSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
	QueryByProjectCustomServicesCustomSkus(context.Context, *model.ProjectDataConfig) ([]model.ProjectCost, error)
}
