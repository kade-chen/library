package impl

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/usagedate"
	"github.com/kade-chen/google-billing-console/apps/usagedate/impl/project"
	"github.com/kade-chen/google-billing-console/apps/usagedate/impl/services"
	"github.com/kade-chen/google-billing-console/apps/usagedate/impl/sku"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

var _ usagedate.ProjectService = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	// col *mongo.Collection
	// token.UnimplementedRPCServer
	ioc.ObjectImpl
	log  *zerolog.Logger
	bq   *bigquery.Client
	svcs usagedate.Service
	skus usagedate.SkuService

	// policy  policy.Service
	// ns      namespace.Service
	// checker security.Checker
	// domain  domain.Service
	// notify  notify.Service
}

func (s *service) Init() error {
	s.log = log.Sub(s.Name())
	s.bq = ioc.Config().Get(configs.AppName).(*impl.Service).BQ

	s.svcs = ioc.Controller().Get(services.AppName).(usagedate.Service)
	s.skus = ioc.Controller().Get(sku.AppName).(usagedate.SkuService)
	return nil
}

func (service) Name() string {
	return project.AppName
}

func (s *service) Close(ctx context.Context) error {
	defer func() {
		if err := s.bq.Close(); err != nil {
			s.log.Error().Msgf("❌ Failed to close BigQuery client: %v", err)
		} else {
			fmt.Println("✅ BigQuery client closed successfully")
		}
	}()
	// 关闭后测试调用
	// q := s.bq.Query("SELECT 1")
	// a, err := q.Run(context.Background())
	// if err != nil {
	// 	fmt.Printf("✅ Client closed: further operations fail as expected: %v\n", err)
	// } else {
	// 	fmt.Println("⚠️ Unexpected: client still appears functional (likely cached connection)")
	// }
	// _ = a
	return nil
}

func (i *service) Priority() int {
	return 0
}
