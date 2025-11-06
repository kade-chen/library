package impl

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/sku"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByDateProjectSKUsAll(ctx context.Context, query string) ([]sku.SkusList, error) {
	startDate := "2025-10-01"
	endDate := "2025-10-01"
	// projectIDs := []string{} // 空数组表示查询全部项目
	projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateProjectSKUsSQL(projectIDs)
	q := s.bq.Query(sql)

	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: startDate},
		{Name: "end_date", Value: endDate},
	}
	if len(projectIDs) > 0 {
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: projectIDs})
	}
	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	var results []sku.SkusList
	// 遍历结果
	for {
		// var r Result
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error reading result: %v", err)
		}

		results = append(results, sku.SkusList{
			ServiceID:      fmt.Sprintf("%v", row["service_id"]),
			SKUId:          fmt.Sprintf("%v", row["sku_id"]),
			SKUDesc:        fmt.Sprintf("%v", row["sku_describe"]),
			ServicePath:    fmt.Sprintf("%v", row["service_path"]),
			ServiceSKUPath: fmt.Sprintf("%v", row["service_sku_path"]),
		})
		// fmt.Printf("| %s | %s | %s | %s | %s\n",
		// 	// formatAny(row["project_id"]),
		// 	format.FormatAny(row["service_id"]),
		// 	format.FormatAny(row["service_description"]),
		// 	format.FormatAny(row["sku_id"]),
		// 	format.FormatAny(row["sku_describe"]),
		// 	format.FormatAny(row["total_rows"]),
		// )
	}

	return results, err
}
