package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	tools"github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataRequest) ([]model.ServicesList, error) {
	// 构造查询
	sql := s.queryByDateProjectSUSQL()
	q := s.bq.Query(sql)
	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)

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
		s.log.Error().Msgf("Failed to execute query, ERROR: %v", err)
		return nil, exception.NewInternalServerError("Failed to execute query, ERROR: %v", err)
	}

	var results []model.ServicesList
	// 遍历结果
	for {
		// var r Result
		var row model.ServicesList
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query finished")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		results = append(results, row)
	}

	return results, err
}
