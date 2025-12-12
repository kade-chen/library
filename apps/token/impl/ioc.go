package impl

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/rs/zerolog"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"

	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	"github.com/kade-chen/library/ioc/config/log"
)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	bq_client   *bigquery.Client
	bq_table    *bigquery.Table
	bqTableFull string
	token.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger

	// policy  policy.Service
	// ns      namespace.Service
	// checker security.Checker
	// domain  domain.Service
	// notify  notify.Service
}

func (s *service) Init() error {
	s.log = log.Sub(s.Name())
	s.bq_client = ioc.Config().Get(configs.AppName).(*impl.Service).BQ
	err := s.bqInit(context.Background())
	if err != nil {
		return err
	}

	// s.ns = ioc.Controller().Get(namespace.AppName).(namespace.Service)
	// s.policy = ioc.Controller().Get(policy.AppName).(policy.Service)
	// s.domain = ioc.Controller().Get(domain.AppName).(domain.Service)
	// s.notify = ioc.Controller().Get(notify.AppName).(notify.Service)

	// s.checker, err = security.NewChecker()
	// if err != nil {
	// 	return fmt.Errorf("new checker error, %s", err)
	// }

	// 初始化所有的auth provider
	if err := provider.Init(); err != nil {
		return err
	}

	token.RegisterRPCServer(grpc.Get().Server(), s)
	projectID := ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID
	dataset := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset
	table := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableToken
	s.bqTableFull = fmt.Sprintf("`%s.%s.%s`", projectID, dataset, table)
	return nil
}

func (service) Name() string {
	return token.AppName
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
	schema, err := bigquery.InferSchema(token.Token{})
	// 使用示例
	makeSchemaNullable(schema)                         // 先把所有字段置 NULLABLE
	setFieldRequired(schema, []string{"access_token"}) // 指定字段为 REQUIRED
	// setFieldRequired(schema, []string{"id", "meta.create_at"}) // 指定字段为 REQUIRED
	// 把 spec.labels 改成 JSON 字段
	forceJSONField(schema, "meta")
	if err != nil {
		s.log.Error().Msgf("infer schema failed, ERROR: %v", err)
		return exception.NewIocRegisterFailed("infer schema failed, ERROR: %v", err)
	}

	table := dataset.Table(ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableToken)

	// ---- 2. 判断 table 是否存在 ----
	_, err = table.Metadata(ctx)
	if err != nil {
		if err == iterator.Done || s.datasetisNotFound(err) {
			s.log.Warn().Msgf("Table: %v not found, creating...", ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableToken)

			// ---- 3. 创建 table ----
			err = table.Create(ctx, &bigquery.TableMetadata{
				Schema: schema,
				// ★★★ 新增：定义主键约束 ★★★
				TableConstraints: &bigquery.TableConstraints{
					PrimaryKey: &bigquery.PrimaryKey{
						// 这里填写 BigQuery 中的列名 (即 struct tag 中的名字)
						Columns: []string{"access_token"},
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

// map 特殊处理转成json，key-value
func forceJSONField(schema bigquery.Schema, fullPath string) {
	parts := strings.Split(fullPath, ".")

	for i := range schema {
		if schema[i].Name == parts[0] {
			// 最后一级 → 当前字段改 JSON
			if len(parts) == 1 {
				schema[i].Type = bigquery.JSONFieldType
				schema[i].Schema = nil // JSON 不能再有 schema
				return
			}

			// 递归下一层
			if schema[i].Type == bigquery.RecordFieldType {
				forceJSONField(schema[i].Schema, strings.Join(parts[1:], "."))
			}
		}
	}
}
