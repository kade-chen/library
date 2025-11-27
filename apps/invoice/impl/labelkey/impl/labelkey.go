package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByInvoiceMonthProjectLabelKeyAll(ctx context.Context, config *model.InvoiceMonthProjectLabelKeyRequest) (model.InvoiceMonthProjectLabelKeyLists, error) {
	// config.StartDate = "2025-10-01"
	// config.EndDate = "2025-10-02"
	// // projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateProjectLabelKeySQL()
	q := s.bq.Query(sql)
	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	s.log.Debug().Msgf("partitionStartTime:%s partitionEndTime:%s", partitionStartTime, partitionEndTime)
	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
	}
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}
	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		return model.InvoiceMonthProjectLabelKeyLists{}, exception.NewInternalServerError("Failed to execute query, ERROR: %v", err)
	}

	var results []model.InvoiceMonthProjectLabelKeyList
	// 遍历结果
	for {
		// var r Result
		var row model.InvoiceMonthProjectLabelKeyList
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query finished")
			break
		}
		if err != nil {
			return model.InvoiceMonthProjectLabelKeyLists{}, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		results = append(results, row)
	}

	return MergeLabelKeyLists(results), err
}

func MergeLabelKeyLists(list []model.InvoiceMonthProjectLabelKeyList) model.InvoiceMonthProjectLabelKeyLists {
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
	return result
}
