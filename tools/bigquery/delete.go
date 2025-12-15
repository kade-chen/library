package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/library/exception"
)

func DeleteSQL(ctx context.Context, sql string, bq_client *bigquery.Client, params []bigquery.QueryParameter) (int64, error) {

	q := bq_client.Query(sql)
	q.Parameters = params

	job, err := q.Run(ctx)
	if err != nil {
		return 0, exception.NewInternalServerError("bq run error: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return 0, exception.NewInternalServerError("bq wait error: %v", err)
	}

	if status.Err() != nil {
		return 0, exception.NewInternalServerError("bq job error: %v", status.Err())
	}
	// ⭐ 正确读取 DML 影响行数
	stats := status.Statistics
	if stats == nil {
		return 0, nil
	}

	qs, ok := stats.Details.(*bigquery.QueryStatistics)
	if !ok || qs.DMLStats == nil {
		// 不是 Query / 不是 DML
		return 0, nil
	}

	deleted := qs.DMLStats.DeletedRowCount
	updated := qs.DMLStats.UpdatedRowCount
	inserted := qs.DMLStats.InsertedRowCount

	affected := deleted + updated + inserted
	return affected, nil
}
