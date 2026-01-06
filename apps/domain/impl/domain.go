package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/domain"
	tools "github.com/kade-chen/google-billing-console/tools/bigquery"
	"github.com/kade-chen/library/exception"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

func (s *service) CreateDomain(ctx context.Context, req *domain.CreateDomainRequest) (*domain.Domain, error) {
	domain, err := domain.NewDomain(req)
	if err != nil {
		return nil, err
	}
	// 查询 SQL
	queryStr := fmt.Sprintf(`SELECT id FROM %s WHERE id = '%s'`, s.bqTableFull, domain.Id)

	bq := s.bq_client.Query(queryStr)

	// 执行 SELECT 并返回结果迭代器
	it, err := bq.Read(ctx)
	if err != nil {
		return domain, exception.NewInternalServerError("query read error: %v", err)
	}

	// 读取一行（你的场景应该只返回 1 行）
	switch err := it.Next(domain); err {
	case iterator.Done:
		s.log.Warn().Msg("domain not exist")
		s.log.Info().Msgf("create %v domain......", domain.Spec.Name)
		// 创建domain
		// 根据结构体自动创建row
		err := tools.BigQueryStructInsert(ctx, s.bq_table, domain)
		if err != nil {
			return nil, err
		}
		s.log.Info().Msgf("create %v domain successful", domain.Spec.Name)
		// 没有查到
		return domain, nil
	case nil:
		s.log.Info().Msgf("domain %v already exist, Do not create duplicates.", domain.Spec.Name)
		return nil, nil
	default:
		s.log.Error().Msgf("iterator error: %v", err)
		return domain, exception.NewInternalServerError("iterator error: %v", err)
	}
}

func (s *service) DescribeDomain(ctx context.Context, req *domain.DescribeDomainRequest) (*domain.Domain, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 构造动态 SQL
	var (
		queryStr string
		params   []bigquery.QueryParameter
	)

	switch req.DescribeBy {
	case domain.DESCRIBE_BY_ID:
		queryStr = fmt.Sprintf(`SELECT * FROM %s WHERE id = @id LIMIT 1`, s.bqTableFull)
		params = []bigquery.QueryParameter{
			{Name: "id", Value: req.Id},
		}
	case domain.DESCRIBE_BY_NAME:
		queryStr = fmt.Sprintf(`SELECT * FROM %s WHERE spec.name = @name LIMIT 1`, s.bqTableFull)
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
		s.log.Error().Msgf("domain: %v not exist", req.Name)
		return nil, exception.NewNotFound("domain: %v not exist", req.Name)
	case nil:
		// 成功查询，row 已是完整 domain
		s.log.Info().Msg("domain query successful")
		return row, nil
	default:
		return nil, exception.NewInternalServerError("iterator error: %v", err)
	}
}

func (s *service) ListDoamin(ctx context.Context, req *domain.ListDomainRequest) (*domain.DomainSet, error) {
	r := domain.NewListDomainRequest(req)
	whereSQL, whereParams := r.WhereSQL()
	pageSQL, pageParams := r.PageSQL()
	params := append(whereParams, pageParams...)

	sql := fmt.Sprintf(`SELECT * FROM %s %s ORDER BY meta.create_at ASC %s`, s.bqTableFull, whereSQL, pageSQL)
	set := domain.NewDomainSet()

	g, ctx := errgroup.WithContext(ctx)
	// --------------------------
	// 1. 并发查询列表
	// --------------------------
	g.Go(func() error {
		var rowCount int64
		q := s.bq_client.Query(sql)
		q.Parameters = params

		it, err := q.Read(ctx)
		if err != nil {
			s.log.Error().Msgf("query domain list error: %v", err)
			return exception.NewInternalServerError("query domain error: %v", err)
		}

		for {
			rowMap := make(map[string]bigquery.Value)
			err := it.Next(&rowMap)

			if err == iterator.Done {
				s.log.Info().Msgf("query domain list done")
				break
			}
			if err != nil {
				s.log.Error().Msgf("query domain list error: %v", err)
				return exception.NewInternalServerError("query domain list error: %v", err)
			}
			// ---- 转 user.User ----
			rowBytes, err := json.Marshal(rowMap)
			if err != nil {
				s.log.Error().Msgf("marshal map error: %v", err)
				return exception.NewInternalServerError("marshal map error: %v", err)
			}
			row := &domain.Domain{}
			if err := json.Unmarshal(rowBytes, row); err != nil {
				s.log.Error().Msgf("unmarshal to domain error: %v", err)
				return exception.NewInternalServerError("unmarshal to domain error: %v", err)
			}
			// row.Desensitization()
			set.Add(row)
			rowCount++
		}
		set.Total = rowCount
		return nil
	})

	// --------------------------
	// 等待两个 goroutine 完成
	// --------------------------
	if err := g.Wait(); err != nil {
		s.log.Error().Msgf("query domain error: %v", err)
		return nil, exception.NewInternalServerError("%v", err)
	}

	return set, nil
}
