package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	logs "github.com/kade-chen/library/ioc/config/log"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"

	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/rs/zerolog"
)

var _ domain.Service = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	bq_client *bigquery.Client
	bq_table  *bigquery.Table
	domain.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (*service) Name() string {
	return domain.AppName
}

func (i *service) Priority() int {
	return 100
}

func (s *service) Init() error {

	s.log = logs.Sub(domain.AppName)
	s.bq_client = ioc.Config().Get(configs.AppName).(*impl.Service).BQ
	err := s.bqInit(context.Background())
	if err != nil {
		return err
	}

	domain.RegisterRPCServer(grpc.Get().Server(), s)
	return err
}

func (s *service) bqInit(ctx context.Context) error {
	dataset := s.bq_client.Dataset(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset)
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
	schema, err := bigquery.InferSchema(domain.Domain{})
	// 使用示例
	makeSchemaNullable(schema)               // 先把所有字段置 NULLABLE
	setFieldRequired(schema, []string{"id"}) // 指定字段为 REQUIRED
	// setFieldRequired(schema, []string{"id", "meta.create_at"}) // 指定字段为 REQUIRED

	if err != nil {
		s.log.Error().Msgf("infer schema failed, ERROR: %v", err)
		return exception.NewIocRegisterFailed("infer schema failed, ERROR: %v", err)
	}

	table := dataset.Table(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableDomain)

	// ---- 2. 判断 table 是否存在 ----
	_, err = table.Metadata(ctx)
	if err != nil {
		if err == iterator.Done || s.datasetisNotFound(err) {
			s.log.Warn().Msgf("Table: %v not found, creating...", ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableDomain)

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
	s.bq_table = table
	s.log.Info().Msg("Table already exists, continue...")
	return nil
}

func (s *service) datasetisNotFound(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		return gerr.Code == 404
	}
	return false
}

// 针对特定字段设置 NULLABLE
func makeSchemaNullable(schema bigquery.Schema) {
	for _, field := range schema {
		field.Required = false // 设置为 NULLABLE
		if field.Type == bigquery.RecordFieldType && field.Schema != nil {
			makeSchemaNullable(field.Schema) // 递归嵌套字段
		}
	}
}

// 针对特定字段设置 REQUIRED
func setFieldRequired(schema bigquery.Schema, fieldNames []string) {
	for _, field := range schema {
		for _, name := range fieldNames {
			if field.Name == name {
				field.Required = true
			}
		}
		if field.Type == bigquery.RecordFieldType && field.Schema != nil {
			setFieldRequired(field.Schema, fieldNames) // 递归嵌套 STRUCT
		}
	}
}

//REPEATED 切片可以自动推断
