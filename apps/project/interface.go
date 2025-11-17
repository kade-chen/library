package project

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "project"
)

type Service interface {
	//辅助功能
	QueryByDateProjectAllServicesAllSkus(context.Context, *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error)
	//by date project
	QueryByDateProject(context.Context, *model.ProjectDataConfig) ([]model.ProjectDateCost, error)

	//by project
	QueryByProject(context.Context, *model.ProjectConfig) ([]model.ProjectCost, error)
}
