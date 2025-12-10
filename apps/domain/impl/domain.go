package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"google.golang.org/api/iterator"
)

func (s *service) CreateDomain(ctx context.Context, req *domain.CreateDomainRequest) (*domain.Domain, error) {
	d, err := domain.NewDomain(req)
	if err != nil {
		return nil, err
	}

	projectID := ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID
	dataset := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset
	table := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableDomain
	tableFull := fmt.Sprintf("`%s.%s.%s`", projectID, dataset, table)
	// 查询 SQL
	queryStr := fmt.Sprintf(`
			SELECT id FROM %s WHERE id = '%s'
			`, tableFull,
		d.Id)

	bq := s.bq_client.Query(queryStr)

	// 执行 SELECT 并返回结果迭代器
	it, err := bq.Read(ctx)
	if err != nil {
		return d, exception.NewInternalServerError("query read error: %v", err)
	}

	// 读取一行（你的场景应该只返回 1 行）
	switch err := it.Next(d); err {
	case iterator.Done:
		s.log.Warn().Msg("domain not exist")
		s.log.Info().Msgf("create %v domain......", d.Spec.Name)
		// 将结构体 d 序列化为 JSON 字节。
		// BigQuery 的 JSON Load Job 只接受 JSON Lines 格式：一条记录一行。
		// 因此需要在序列化后的 JSON 后面追加换行符 '\n'。
		b, err := json.Marshal(d)
		if err != nil {
			return nil, nil
		}
		data := append(b, '\n')

		// 使用 bytes.Reader 包装内存 JSON 数据，作为 BigQuery Load Job 的输入源。
		// 这是不经过 Streaming Insert、也不需要上传到 GCS 的“内存加载方式”。
		rdr := bytes.NewReader(data)

		// 创建一个 ReaderSource，告知 BigQuery 数据来源是一个 Reader。
		// 默认格式是 CSV，因此必须指定为 JSON，否则会报解析错误。
		src := bigquery.NewReaderSource(rdr)
		src.SourceFormat = bigquery.JSON

		// 基于 src 创建一个 Loader（加载任务），
		// 这是执行 BigQuery Load Job 的核心对象。
		// 它不会像 Streaming Insert 一样写入 Streaming Buffer，因此可避免 UPDATE/DELETE 受限的问题。
		inserter := s.bq_table.LoaderFrom(src)

		// 设置写入模式（Write Disposition）。
		// WriteAppend 表示追加写入：不会清空表，不会覆盖数据，只是往表尾追加新记录。
		// 其他模式：WriteTruncate 会全部清空再写入；WriteEmpty 表为空时才能写。
		inserter.WriteDisposition = bigquery.WriteAppend

		job, err := inserter.Run(ctx)
		if err != nil {
			return nil, exception.NewInternalServerError("load job run error: %v", err)
		}

		status, err := job.Wait(ctx)
		if err != nil {
			return nil, exception.NewInternalServerError("wait error: %v", err)
		}

		if status.Err() != nil {
			// 打印 BigQuery loader 的详细错误
			for _, e := range status.Errors {
				s.log.Error().Msgf("BQ Load Error: %+v", e)
			}
			return nil, exception.NewInternalServerError("load job status error: %v", status.Err())
		}
		s.log.Info().Msgf("create %v domain successful", d.Spec.Name)
		// 没有查到
		return d, nil
	case nil:
		s.log.Info().Msgf("domain %v already exist, query successful", d.Spec.Name)
		return d, nil
	default:
		s.log.Error().Msgf("iterator error: %v", err)
		return d, exception.NewInternalServerError("iterator error: %v", err)
	}
}

func (s *service) DescribeDomain(ctx context.Context, req *domain.DescribeDomainRequest) (*domain.Domain, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 表名构造
	projectID := ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID
	dataset := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset
	table := ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableDomain
	tableFull := fmt.Sprintf("`%s.%s.%s`", projectID, dataset, table)

	// 构造动态 SQL
	var (
		queryStr string
		params   []bigquery.QueryParameter
	)

	switch req.DescribeBy {
	case domain.DESCRIBE_BY_ID:
		queryStr = fmt.Sprintf(`
			SELECT * FROM %s
			WHERE id = @id
			LIMIT 1
		`, tableFull)
		params = []bigquery.QueryParameter{
			{Name: "id", Value: req.Id},
		}

	case domain.DESCRIBE_BY_NAME:
		queryStr = fmt.Sprintf(`
			SELECT * FROM %s
			WHERE spec.name = @name
			LIMIT 1
		`, tableFull)
		params = []bigquery.QueryParameter{
			{Name: "name", Value: req.Name},
		}

	default:
		return nil, exception.NewBadRequest("invalid DescribeBy type")
	}

	// 执行查询
	q := s.bq_client.Query(queryStr)
	q.Parameters = params

	it, err := q.Read(ctx)
	if err != nil {
		return nil, exception.NewInternalServerError("query error: %v", err)
	}

	// ⚠ 注意：row 必须完整匹配 BigQuery 表字段
	// 你可以直接用 domain.Domain（推荐）
	row := &domain.Domain{}

	err = it.Next(row)
	switch err {
	case iterator.Done:
		s.log.Error().Msg("domain not exist")
		return nil, exception.NewNotFound("domain not exist")
	case nil:
		// 成功查询，row 已是完整 domain
		s.log.Info().Msg("domain query successful")
		return row, nil
	default:
		return nil, exception.NewInternalServerError("iterator error: %v", err)
	}
}
