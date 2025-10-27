package impl

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByDateProject(ctx context.Context, query string) error {
	startDate := "2025-10-01"
	endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateProjectSQL(projectIDs)
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
		fmt.Printf("%s | %s | %s | %s | %s | %s\n",
			formatAny(row["usage_date"]),
			formatAny(row["project_id"]),
			formatAny(row["cost_list"]),
			formatAny(row["cost_list_abs"]),
			formatAny(row["cost"]),
			formatAny(row["total_cost_sum"]),
		)
	}
	return nil
}

func formatAny(n any) any {
	switch v := n.(type) {
	case bigquery.NullFloat64:
		if v.Valid {
			return fmt.Sprintf("%.2f", v.Float64) // 保留两位小数
		}
	case bigquery.NullString:
		if v.Valid {
			return v.String() // 注意大括号：调用方法
		}
	case bigquery.NullDate:
		if v.Valid {
			return v.Date // 转成 yyyy-mm-dd
		}
	case float64:
		return fmt.Sprintf("%.2f", v)
	case int64:
		return fmt.Sprintf("%2d", v)
	case string:
		return v
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
	return n
}
