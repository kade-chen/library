package billing

import "context"

const (
	AppName = "bigquery"
)

type Service interface {
	// QueryByProject(ctx context.Context, query string) error
	QueryByDateProject(ctx context.Context, query string) error
	// QueryBySku(ctx context.Context, query string) error
	// QueryByDateSku(ctx context.Context, query string) error
	// QueryByService(ctx context.Context, query string) error
	// QueryByDateService(ctx context.Context, query string) error
}
