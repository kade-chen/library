package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/google-billing-console/tools/trances"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByInvoiceMonthProjectLabelKeyAll(ctx context.Context, config *model.InvoiceMonthProjectLabelKeyRequest) (model.InvoiceMonthProjectLabelKeyLists, error) {
	if ctx.Value("trances_id") == nil {
		ctx = context.WithValue(context.Background(), "trances_id", trances.NewTraceID())
	}
	s.log.Info().Msgf("trances_id=%s, The User begins Query for UsageDateByDatePojectAPI", ctx.Value("trances_id"))
	// 构造查询
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL......", ctx.Value("trances_id"))
	sql := s.queryByDateProjectLabelKeySQL()
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL successful", ctx.Value("trances_id"))

	q := s.bq.Query(sql)
	s.log.Info().Msgf("trances_id=%s, Configuring query parameters......", ctx.Value("trances_id"))

	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	s.log.Info().Msgf("trances_id=%s, start_date=%v, end_date=%v, PartitionStartTime=%v, PartitionEndTime=%v", ctx.Value("trances_id"), config.StartDate, config.EndDate, partitionStartTime, partitionEndTime)

	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
	}
	s.log.Info().Msgf("trances_id=%s, Retrieving project_ids......", ctx.Value("trances_id"))
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, project_ids=%v", ctx.Value("trances_id"), config.ProjectIDs)
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, project_ids is empty, query all project_ids", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}
	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
		return model.InvoiceMonthProjectLabelKeyLists{}, exception.NewInternalServerError("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
	}

	var results []model.InvoiceMonthProjectLabelKeyList
	// 遍历结果
	for {
		// var r Result
		var row model.InvoiceMonthProjectLabelKeyList
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msgf("trances_id=%s, Query completed for UsageDateByDatePojectAPI", ctx.Value("trances_id"))
			break
		}
		if err != nil {
			s.log.Error().Msgf("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
			return model.InvoiceMonthProjectLabelKeyLists{}, exception.NewInternalServerError("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
		}
		results = append(results, row)
	}

	return s.mergeLabelKeyLists(ctx, results), err
}

func (s *service) mergeLabelKeyLists(ctx context.Context, list []model.InvoiceMonthProjectLabelKeyList) model.InvoiceMonthProjectLabelKeyLists {
	result := model.InvoiceMonthProjectLabelKeyLists{
		LabelKey:     make([]string, 0),
		LabelValue:   make([]string, 0),
		LabelKeyPath: make([]string, 0),
	}
	// 用 map 做去重
	keyMap := make(map[string]struct{})
	valueMap := make(map[string]struct{})
	pathMap := make(map[string]struct{})

	for _, item := range list {
		// LabelKey 去重
		if item.LabelKey.Valid {
			if _, exists := keyMap[item.LabelKey.String()]; !exists {
				keyMap[item.LabelKey.String()] = struct{}{}
				result.LabelKey = append(result.LabelKey, item.LabelKey.String())
			}
		}

		// LabelValue 去重
		if item.LabelValue.Valid {
			if _, exists := valueMap[item.LabelValue.String()]; !exists {
				valueMap[item.LabelValue.String()] = struct{}{}
				result.LabelValue = append(result.LabelValue, item.LabelValue.String())
			}
		}

		// LabelKeyPath 去重
		if item.LabelKeyPath.Valid {
			if _, exists := pathMap[item.LabelKeyPath.String()]; !exists {
				pathMap[item.LabelKeyPath.String()] = struct{}{}
				result.LabelKeyPath = append(result.LabelKeyPath, item.LabelKeyPath.String())
			}
		}
	}
	s.log.Info().Msgf("trances_id=%s, LabelKey is Deduplication", ctx.Value("trances_id"))
	s.log.Info().Msgf("trances_id=%s, LabelValue is Deduplication", ctx.Value("trances_id"))
	s.log.Info().Msgf("trances_id=%s, LabelKeyPath is Deduplication", ctx.Value("trances_id"))
	return result
}
