package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/token"
	tools "github.com/kade-chen/google-billing-console/tools/bigquery"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"
	// "go.mongodb.org/mongo-driver/mongo"
)

// query wether  the token  exists for mongdb
func (s *service) get(ctx context.Context, id string) (*token.Token, error) {

	sql := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE access_token = @access_token
		LIMIT 1
	`, s.bqTableFull)

	// 执行查询
	q := s.bq_client.Query(sql)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "access_token", Value: id},
	}

	it, err := q.Read(ctx)
	if err != nil {
		return nil, exception.NewInternalServerError("query error: %v", err)
	}

	// --- 修改开始 ---
	// 1. 使用通用 Map 接收 BigQuery 数据
	// BigQuery 客户端可以很好地将整行数据映射到 map[string]bigquery.Value
	rowMap := make(map[string]bigquery.Value)
	err = it.Next(&rowMap)
	if err == iterator.Done {
		return nil, exception.NewNotFound("token not exist")
	} else if err != nil {
		return nil, err
	}

	// --- 修复 Start ---
	// BigQuery STRUCT → Go map, 不需要额外处理，只要确保不是 string
	if metaRaw, ok := rowMap["meta"]; ok {
		switch v := metaRaw.(type) {
		case string:
			// 如果 BigQuery 错误地把 STRUCT 当 JSON string 返回
			var m map[string]interface{}
			if json.Unmarshal([]byte(v), &m) == nil {
				rowMap["meta"] = m
			}
		case map[string]any, map[string]bigquery.Value:
			// ok 不处理
		case nil:
			// ok
		default:
			// 不支持其他类型，置空
			rowMap["meta"] = nil
		}
	}
	// --- 修复 End ---

	rowBytes, err := json.Marshal(rowMap)
	if err != nil {
		return nil, err
	}

	row := &token.Token{}
	if err := protojson.Unmarshal(rowBytes, row); err != nil {
		return nil, fmt.Errorf("protojson unmarshal error: %w", err)
	}

	return row, nil
}

func (s *service) save(ctx context.Context, tk *token.Token) error {
	// 创建用户
	// 根据结构体自动创建row
	err := tools.BigQueryStructInsert(ctx, s.bq_table, tk)
	if err != nil {
		s.log.Error().Msgf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

func (s *service) delete(ctx context.Context, ins *token.Token) error {
	// if ins == nil || ins.AccessToken == "" {
	// 	return fmt.Errorf("access tpken is nil")
	// }

	// result, err := s.col.DeleteOne(ctx, bson.M{"_id": ins.AccessToken})
	// if err != nil {
	// 	return exception.NewInternalServerError("delete token(%s) error, %s", ins.AccessToken, err)
	// }

	// if result.DeletedCount == 0 {
	// 	return exception.NewNotFound("book %s not found", ins.AccessToken)
	// }

	return nil
}
