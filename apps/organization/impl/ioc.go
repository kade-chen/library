package impl

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	logs "github.com/kade-chen/library/ioc/config/log"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"

	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/organization"
	"github.com/rs/zerolog"
)

var _ organization.Service = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	bq_client   *bigquery.Client
	bq_table    *bigquery.Table
	bqTableFull string
	organization.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (*service) Name() string {
	return organization.AppName
}

func (i *service) Priority() int {
	return 100
}

func (s *service) Init() error {

	s.log = logs.Sub(organization.AppName)
	s.bq_client = ioc.Config().Get(configs.AppName).(*impl.Service).BQ
	err := s.bqInit(context.Background())
	if err != nil {
		return err
	}

	organization.RegisterRPCServer(grpc.Get().Server(), s)

	projectID := ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID
	dataset := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset
	table := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableOrganization
	s.bqTableFull = fmt.Sprintf("`%s.%s.%s`", projectID, dataset, table)

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
	schema, err := bigquery.InferSchema(organization.Organization{})
	// 使用示例
	makeSchemaNullable(schema)               // 先把所有字段置 NULLABLE
	setFieldRequired(schema, []string{"id"}) // 指定字段为 REQUIRED
	// setFieldRequired(schema, []string{"id", "meta.create_at"}) // 指定字段为 REQUIRED

	if err != nil {
		s.log.Error().Msgf("infer schema failed, ERROR: %v", err)
		return exception.NewIocRegisterFailed("infer schema failed, ERROR: %v", err)
	}

	table := dataset.Table(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableOrganization)

	// ---- 2. 判断 table 是否存在 ----
	_, err = table.Metadata(ctx)
	if err != nil {
		if err == iterator.Done || s.datasetisNotFound(err) {
			s.log.Warn().Msgf("Table: %v not found, creating...", ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableOrganization)

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
	for i := range schema {
		schema[i].Required = false
		if schema[i].Type == bigquery.RecordFieldType && schema[i].Schema != nil {
			makeSchemaNullable(schema[i].Schema)
		}
	}
}

// 针对特定字段设置 REQUIRED
func setFieldRequired(schema bigquery.Schema, fieldNames []string) {
	for i := range schema {
		for _, name := range fieldNames {
			if schema[i].Name == name {
				schema[i].Required = true
			}
		}
		if schema[i].Type == bigquery.RecordFieldType && schema[i].Schema != nil {
			setFieldRequired(schema[i].Schema, fieldNames)
		}
	}
}

//REPEATED 切片可以自动推断
