package sku

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "skus"
)

type Service interface {
	QueryByDateProjectSKUsAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.SkusList, error)

	//by date sku
	QueryByDateSku(context.Context, *model.SkuDataConfig) ([]model.SkuDateCost, error)

	// //by Service
	QueryBySku(context.Context, *model.SkuConfig) ([]model.SkuCost, error)
}
