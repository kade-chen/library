package sku

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/project"
)

const (
	AppName = "skus"
)

type Service interface {
	QueryByDateProjectSKUsAll(ctx context.Context, config *project.ProjectDataConfig) ([]model.SkusList, error)
}
