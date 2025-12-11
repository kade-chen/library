package bigquery

import (
	"bytes"
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/library/exception"
)

// 根据结构体自动创建row
func BigQueryStructInsert(ctx context.Context, bq_table *bigquery.Table, any any) error {
	b, err := json.Marshal(any)
	if err != nil {
		return nil
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
	inserter := bq_table.LoaderFrom(src)

	// 设置写入模式（Write Disposition）。
	// WriteAppend 表示追加写入：不会清空表，不会覆盖数据，只是往表尾追加新记录。
	// 其他模式：WriteTruncate 会全部清空再写入；WriteEmpty 表为空时才能写。
	inserter.WriteDisposition = bigquery.WriteAppend

	job, err := inserter.Run(ctx)
	if err != nil {
		return exception.NewInternalServerError("load job run error: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return exception.NewInternalServerError("wait error: %v", err)
	}

	if status.Err() != nil {
		// 打印 BigQuery loader 的详细错误
		// for _, e := range status.Errors {
		// 	s.log.Error().Msgf("BQ Load Error: %+v", e)
		// }
		return exception.NewInternalServerError("load job status error: %v", status.Err())
	}
	return nil
}
