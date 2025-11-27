package impl

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/labelkey"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
)

var _ invoice.LabelKeyService = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	// col *mongo.Collection
	// token.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger
	bq  *bigquery.Client

	// policy  policy.Service
	// ns      namespace.Service
	// checker security.Checker
	// domain  domain.Service
	// notify  notify.Service
}

func (s *service) Init() error {
	s.log = log.Sub(s.Name())
	client, err := bigquery.NewClient(context.Background(), ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID, option.WithCredentialsFile(ioc.Config().Get(configs.AppName).(*impl.Service).Default_Service_Account_Name))
	if err != nil {
		s.log.Error().Msgf("Failed to create BigQuery client: %v", err)
		return exception.NewIocRegisterFailed("Failed to create BigQuery client: %v", err)
	}
	s.bq = client
	// 验证能否列出 dataset
	it := client.Datasets(context.Background())
	dataset, err := it.Next()
	if err == iterator.Done {
		s.log.Debug().Msg("⚠️ No datasets found, but client works fine.")
	} else if err != nil {
		s.log.Error().Msgf("❌ Failed to verify connection: %v", err)
		return exception.NewIocRegisterFailed("❌ Failed to verify connection: %v", err)
	} else {
		s.log.Info().Msgf("✅ Verified connection! Example dataset: %s", dataset.DatasetID)
	}
	return nil
}

func (service) Name() string {
	return labelkey.AppName
}

func (s *service) Close(ctx context.Context) error {
	defer func() {
		if err := s.bq.Close(); err != nil {
			fmt.Printf("❌ Failed to close BigQuery client: %v\n", err)
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
	return 1
}
