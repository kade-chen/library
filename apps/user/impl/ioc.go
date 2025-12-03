package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/rs/zerolog"
	"google.golang.org/api/iterator"

	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	logs "github.com/kade-chen/library/ioc/config/log"
)

var _ user.Service = &service{}

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	ioc.ObjectImpl
	col *bigquery.Client
	log *zerolog.Logger

	user.UnimplementedRPCServer
}

func (i *service) Name() string {
	return user.AppName
}

func (i *service) Priority() int {
	return 11
}

func (s *service) Init() error {
	s.log = logs.Sub(user.AppName)
	s.col = ioc.Config().Get(configs.AppName).(*impl.Service).BQ
	err := s.bqInit(context.Background())
	if err != nil {
		return err
	}
	user.RegisterRPCServer(grpc.Get().Server(), s)

	// conn, _ := grpc1.Dial("localhost:1234", grpc1.WithTransportCredentials(insecure.NewCredentials()))
	// v := user.NewRPCClient(conn)
	// v.
	return nil
}

func (s *service) bqInit(ctx context.Context) error {
	dataset := s.col.Dataset(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset)
	_, err := dataset.Metadata(ctx)
	if err != nil {
		// 如果不存在，则创建
		if s.datasetisNotFound(err) {
			s.log.Warn().Msg("Dataset not found, creating...")
			if err := dataset.Create(ctx, &bigquery.DatasetMetadata{
				Location: "asia-east1", // ★ 必填，一般是 US / EU / asia-east1 / asia-northeast1
			}); err != nil {
				s.log.Error().Msgf("create dataset failed, ERROR: %v", err)
				return exception.NewIocRegisterFailed("create dataset failed, ERROR: %v", err)
			}
			s.log.Info().Msgf("Dataset created! Dataset ID: %v", dataset.DatasetID)
		} else {
			s.log.Error().Msgf("get dataset metadata failed, ERROR: %v", err)
		}
	}

	// ---- 1. 自动从结构体推断 schema ----
	schema, err := bigquery.InferSchema(user.User{})
	if err != nil {
		s.log.Error().Msgf("infer schema failed, ERROR: %v", err)
		return exception.NewIocRegisterFailed("infer schema failed, ERROR: %v", err)
	}
	table := dataset.Table(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableUser)

	// ---- 2. 判断 table 是否存在 ----
	_, err = table.Metadata(ctx)
	if err != nil {
		if err == iterator.Done || s.datasetisNotFound(err) {
			s.log.Warn().Msgf("Table: %v not found, creating...", ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableUser)

			// ---- 3. 创建 table ----
			err = table.Create(ctx, &bigquery.TableMetadata{
				Schema: schema,
				// ★★★ 新增：定义主键约束 ★★★
				TableConstraints: &bigquery.TableConstraints{
					PrimaryKey: &bigquery.PrimaryKey{
						// 这里填写 BigQuery 中的列名 (即 struct tag 中的名字)
						Columns: []string{"id"},
					},
				},
			})
			if err != nil {
				s.log.Error().Msgf("create table failed, ERROR: %v", err)
				return exception.NewIocRegisterFailed("create table failed, ERROR: %v", err)
			}
			s.log.Info().Msg("Table created successfully")
		} else {
			s.log.Error().Msgf("failed to get table metadata, ERROR: %v", err)
			return exception.NewIocRegisterFailed("failed to get table metadata, ERROR: %v", err)
		}
	}
	s.log.Info().Msg("Table already exists, continue...")
	return nil
}
