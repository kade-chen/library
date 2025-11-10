package services

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/common/model"
)

const (
	AppName = "services"
)

type Service interface {
	QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.ServicesList, error)
}

// type ServicesList struct {
// 	ServiceID   bigquery.NullString `bigquery:"service_id" json:"service_id"`
// 	ServiceDesc bigquery.NullString `bigquery:"service_description" json:"service_description"`
// 	ServicePath bigquery.NullString `bigquery:"service_path" json:"service_path"`
// }
