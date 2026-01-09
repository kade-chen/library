package impl

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/golang-jwt/jwt/v5"
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
	BQ                                           *bigquery.Client
	log                                          *zerolog.Logger
	Default_Project_ID                           string          `toml:"default_project_id" json:"default_project_id" yaml:"default_project_id"`
	Default_Service_Account_Name                 string          `toml:"default_service_account_name" json:"default_service_account_name" yaml:"default_service_account_name"`
	GoogleBillingConsoleDataset                  string          `toml:"google_billing_console_dataset" json:"google_billing_console_dataset" yaml:"google_billing_console_dataset"`
	GoogleBillingConsoleDatasetTableUser         string          `toml:"google_billing_console_dataset_table_user" json:"google_billing_console_dataset_table_user" yaml:"google_billing_console_dataset_table_user"`
	GoogleBillingConsoleDatasetTableOrganization string          `toml:"google_billing_console_dataset_table_organization" json:"google_billing_console_dataset_table_organization" yaml:"google_billing_console_dataset_table_organization"`
	GoogleBillingConsoleDatasetTableToken        string          `toml:"google_billing_console_dataset_table_token" json:"google_billing_console_dataset_table_token" yaml:"google_billing_console_dataset_table_token"`
	JwtPublicPemFile                             string          `toml:"jwt_public_pem_file" json:"jwt_public_pem_file" yaml:"jwt_public_pem_file"`
	JwtPrivatePemFile                            string          `toml:"jwt_private_pem_file" json:"jwt_private_pem_file" yaml:"jwt_private_pem_file"`
	JwtPublicKey                                 *rsa.PublicKey  `toml:"-" json:"-" yaml:"-"`
	JwtPrivateKey                                *rsa.PrivateKey `toml:"-" json:"-" yaml:"-"`
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
		s.log.Debug().Msgf("✅ Verified connection! Example dataset: %s", dataset.DatasetID)
	}
	s.BQ = client

	err = s.JWTInit(context.Background())
	if err != nil {
		return err
	}
	s.log.Debug().Msgf("%s init successful", s.Name())
	return nil
}

func (Service) Name() string {
	return configs.AppName
}

func (i *Service) Priority() int {
	return 0
}
func (i *Service) Close(ctx context.Context) error {
	defer func() {
		if err := i.BQ.Close(); err != nil {
			i.log.Error().Msgf("❌ Failed to close BigQuery client: %v", err)
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

func (i *Service) JWTInit(ctx context.Context) error {
	// 读取整个文件为字节切片
	i.log.Debug().Msgf("✅ JWT Certificate Init......")
	i.log.Debug().Msgf("✅ Reading public key file: %v", i.JwtPublicPemFile)
	data, err := os.ReadFile(i.JwtPublicPemFile)
	if err != nil {
		i.log.Error().Msgf("❌ Failed to read public key file: %v", err)
		return exception.NewIocApiRegisterFailed("read public key file failed: %v", err)
	}
	i.log.Debug().Msgf("✅ Public key file read successfully: %v", i.JwtPublicPemFile)

	i.log.Debug().Msgf("✅ Parsing public key......")
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		i.log.Error().Msgf("❌ Failed to parse public key: %v", err)
		return exception.NewIocApiRegisterFailed("parse public key failed: %v", err)
	}
	i.JwtPublicKey = publicKey
	i.log.Debug().Msgf("✅ Public key parsed successfully")

	i.log.Debug().Msgf("✅ Reading private key file: %v", i.JwtPrivatePemFile)
	data, err = os.ReadFile(i.JwtPrivatePemFile)
	if err != nil {
		i.log.Error().Msgf("❌ Failed to read private key file: %v", err)
		return exception.NewIocApiRegisterFailed("read private key file failed: %v", err)
	}
	i.log.Debug().Msgf("✅ Private key file read successfully: %v", i.JwtPrivatePemFile)

	i.log.Debug().Msgf("✅ Parsing private key......")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		i.log.Error().Msgf("❌ Failed to parse private key: %v", err)
		return exception.NewIocApiRegisterFailed("parse private key failed: %v", err)
	}
	i.JwtPrivateKey = privateKey
	i.log.Debug().Msgf("✅ Private key parsed successfully")

	i.log.Debug().Msgf("✅ JWT Certificate initialized successfully")
	return nil
}
