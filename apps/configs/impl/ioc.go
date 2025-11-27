package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// var _ configs.Service = (*Service)(nil)

func init() {
	ioc.Config().Registry(&Service{})
}

type Service struct {
	ioc.ObjectImpl
	BQ                           *bigquery.Client
	log                          *zerolog.Logger
	Default_Project_ID           string `toml:"default_project_id" json:"default_project_id" yaml:"default_project_id"`
	Default_Service_Account_Name string `toml:"default_service_account_name" json:"default_service_account_name" yaml:"default_service_account_name"`
}

func (s *Service) Init() error {
	s.log = log.Sub(s.Name())
	s.log.Debug().Msgf("default_project_id:%s default_service_account_name:%s", s.Default_Project_ID, s.Default_Service_Account_Name)
	s.log.Debug().Msgf("bq client begin init......")
	client, err := bigquery.NewClient(context.Background(), s.Default_Project_ID, option.WithCredentialsFile(s.Default_Service_Account_Name))
	if err != nil {
		s.log.Error().Msgf("Failed to create BigQuery client: %v", err)
		return exception.NewIocRegisterFailed("Failed to create BigQuery client: %v", err)
	}

	// 验证能否列出 dataset
	it := client.Datasets(context.Background())
	dataset, err := it.Next()
	if err == iterator.Done {
		s.log.Debug().Msgf("No datasets found, but client works fine.")
	} else if err != nil {
		s.log.Error().Msgf("❌ Failed to verify connection: %v", err)
		return exception.NewIocRegisterFailed("❌ Failed to verify connection: %v", err)
	} else {
		s.log.Debug().Msgf("✅ Verified connection! Example dataset: %s\n", dataset.DatasetID)
	}
	s.BQ = client
	return nil
}

func (Service) Name() string {
	return configs.AppName
}

func (i *Service) Priority() int {
	return 0
}
