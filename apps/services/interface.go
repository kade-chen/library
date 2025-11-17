package services

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "services"
)

type Service interface {
	//辅助功能
	QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.ServicesList, error)
	//by date Service
	QueryByDateService(context.Context, *model.ServiceDataConfig) ([]model.ServiceDateCost, error)

	// //by Service
	QueryByService(context.Context, *model.ServiceConfig) ([]model.ServiceCost, error)
}
