package impl

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/google-billing-console/apps/tools"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByDateProjectSKUsAll(ctx context.Context, config *project.ProjectDataConfig) ([]model.SkusList, error) {
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateProjectSKUsSQL(projectIDs)
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
		log.Fatalf("Failed to execute query: %v", err)
	}

	var results []model.SkusList
	// 遍历结果
	for {
		// var r Result
		var row model.SkusList
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error reading result: %v", err)
		}
		results = append(results, row)
	}

	return results, err
}
