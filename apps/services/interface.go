package services

import "context"

const (
	AppName = "services"
)

type Service interface {
	QueryByDateProjectServicesAll(ctx context.Context, query string) ([]ServicesList, error)
}

type ServicesList struct {
	ServiceID   string `json:"service_id"`
	ServiceDesc string `json:"service_description"`
	ServicePath string `json:"service_path"`
}
