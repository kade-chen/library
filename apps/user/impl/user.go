package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/user"
	tools "github.com/kade-chen/google-billing-console/tools/bigquery"
	"github.com/kade-chen/library/exception"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

func (s *service) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	user, err := user.NewUser(req)
	if err != nil {
		return nil, err
	}
	//
	// 查询 SQL
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE spec.username="%s" AND spec.domain="%s" LIMIT 1`, s.bqTableFull, user.Spec.Username, user.Spec.Domain)

	bq := s.bq_client.Query(sql)
	// 执行 SELECT 并返回结果迭代器
	it, err := bq.Read(ctx)
	if err != nil {
		return user, exception.NewInternalServerError("query read error: %v", err)
	}
	// 读取一行（你的场景应该只返回 1 行）
	switch err := it.Next(user); err {
	case iterator.Done:
		s.log.Warn().Msg("user not exist")
		s.log.Info().Msgf("create %v user......", user.Spec.Username)
		// 创建用户
		// 根据结构体自动创建row
		err := tools.BigQueryStructInsert(ctx, s.bq_table, user)
		if err != nil {
			return nil, err
		}
		s.log.Info().Msgf("create %v user successful", user.Spec.Username)
		// 没有查到
		return user, nil
	case nil:
		s.log.Info().Msgf("domain %v already exist, Do not create duplicates.", user.Spec.Username)
		return nil, nil
	default:
		s.log.Error().Msgf("iterator error: %v", err)
		return user, exception.NewInternalServerError("iterator error: %v", err)
	}
}

func (s *service) DescribeUser(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	// 构造动态 SQL
	var (
		sql    string
		params []bigquery.QueryParameter
	)

	switch req.DescribeBy {
	case user.DESCRIBE_BY_USER_ID:
		sql = fmt.Sprintf(`SELECT * FROM %s WHERE id = @id LIMIT 1`, s.bqTableFull)
		params = []bigquery.QueryParameter{
			{Name: "id", Value: req.Id},
		}
	case user.DESCRIBE_BY_USER_NAME:
		sql = fmt.Sprintf(`SELECT * FROM %s WHERE spec.username = @name AND spec.domain = @domain LIMIT 1`, s.bqTableFull)
		params = []bigquery.QueryParameter{
			{Name: "name", Value: req.Username},
			{Name: "domain", Value: req.Domain},
		}
	// case user.DESCRIBE_BY_FEISHU_USER_ID:
	// 	queryStr = fmt.Sprintf(`
	// 		SELECT * FROM %s
	// 		WHERE spec.name = @name
	// 		LIMIT 1
	// 	`, s.bqTableFull)
	// 	params = []bigquery.QueryParameter{
	// 		{Name: "name", Value: req.Username},
	// 	}
	default:
		return nil, exception.NewBadRequest("unknow desribe by %s", req.DescribeBy)
	}
	// 执行查询
	q := s.bq_client.Query(sql)
	q.Parameters = params

	it, err := q.Read(ctx)
	if err != nil {
		return nil, exception.NewInternalServerError("query error: %v", err)
	}

	// --- 修改开始 ---
	// 1. 使用通用 Map 接收 BigQuery 数据
	// BigQuery 客户端可以很好地将整行数据映射到 map[string]bigquery.Value
	rowMap := make(map[string]bigquery.Value)

	err = it.Next(&rowMap)
	switch err {
	case iterator.Done:
		s.log.Error().Msg("user not exist")
		return nil, exception.NewNotFound("user not exist")
	case nil:
		// 2. 特殊处理 labels 字段
		// BigQuery 的 JSON 类型在 Go 中会被读作 string，如果不处理直接转 JSON，
		// structpb.Struct 会因为接收到 string 而非 object 报错。
		if spec, ok := rowMap["spec"].(map[string]bigquery.Value); ok {
			if labelsRaw, ok := spec["labels"]; ok && labelsRaw != nil {
				// 如果 BigQuery 返回的是 JSON 字符串，我们需要先解开它
				if jsonStr, ok := labelsRaw.(string); ok {
					var labelsMap map[string]interface{}
					// 尝试解析 JSON 字符串
					if err := json.Unmarshal([]byte(jsonStr), &labelsMap); err == nil {
						// 将字符串替换为解析后的 Map对象
						spec["labels"] = labelsMap
					}
				}
				// 如果 BigQuery 返回的已经是 map (如果是 STRUCT 类型)，则无需处理
			}
		}

		// 3. 将 Map 转换为 user.User 结构体
		// 利用 JSON 作为中转： Map -> JSON Bytes -> User Struct
		// 因为 user.User 生成的代码中有 json:"labels" 标签，这种方法完美兼容
		rowBytes, err := json.Marshal(rowMap)
		if err != nil {
			return nil, exception.NewInternalServerError("marshal map error: %v", err)
		}

		row := &user.User{}
		// 注意：使用标准库 json.Unmarshal，因为它支持 @gotags 生成的 json tag
		// protojson.Unmarshal 可能不支持结构体 Tag 别名
		if err := json.Unmarshal(rowBytes, row); err != nil {
			return nil, exception.NewInternalServerError("unmarshal to user error: %v", err)
		}

		s.log.Info().Msg("user query successful")
		return row, nil
	default:
		return nil, exception.NewInternalServerError("iterator error: %v", err)
	}
	// --- 修改结束 ---

	// //3.initalize null value
	// u.InitializeNullValue()
}

// 查询用户列表

func (s *service) ListUser(ctx context.Context, req *user.QueryUserRequest) (*user.UserSet, error) {
	r := user.NewQueryUserRequest(req)
	whereSQL, whereParams := r.WhereSQL()
	pageSQL, pageParams := r.PageSQL()
	params := append(whereParams, pageParams...)

	sql := fmt.Sprintf(`SELECT * FROM %s %s ORDER BY meta.create_at ASC %s`, s.bqTableFull, whereSQL, pageSQL)
	set := user.NewUserSet()

	g, ctx := errgroup.WithContext(ctx)
	// --------------------------
	// 1. 并发查询列表
	// --------------------------
	g.Go(func() error {
		var rowCount int64
		if !r.SkipItems {
			q := s.bq_client.Query(sql)
			q.Parameters = params

			it, err := q.Read(ctx)
			if err != nil {
				s.log.Error().Msgf("query user list error: %v", err)
				return exception.NewInternalServerError("query user error: %v", err)
			}

			for {
				rowMap := make(map[string]bigquery.Value)
				err := it.Next(&rowMap)

				if err == iterator.Done {
					s.log.Info().Msgf("query user list done")
					break
				}
				if err != nil {
					s.log.Error().Msgf("query user list error: %v", err)
					return exception.NewInternalServerError("query user list error: %v", err)
				}

				// ---- 处理 labels ----
				if spec, ok := rowMap["spec"].(map[string]bigquery.Value); ok {
					if labelsRaw, ok := spec["labels"]; ok && labelsRaw != nil {
						if jsonStr, ok := labelsRaw.(string); ok {
							var labelsMap map[string]interface{}
							if err := json.Unmarshal([]byte(jsonStr), &labelsMap); err == nil {
								spec["labels"] = labelsMap
							}
						}
					}
				}

				// ---- 转 user.User ----
				rowBytes, err := json.Marshal(rowMap)
				if err != nil {
					s.log.Error().Msgf("marshal map error: %v", err)
					return exception.NewInternalServerError("marshal map error: %v", err)
				}

				row := &user.User{}
				if err := json.Unmarshal(rowBytes, row); err != nil {
					s.log.Error().Msgf("unmarshal to user error: %v", err)
					return exception.NewInternalServerError("unmarshal to user error: %v", err)
				}

				row.Desensitization()
				set.Add(row)
				rowCount++
			}
			set.Total = rowCount
		}

		return nil
	})

	// --------------------------
	// 等待两个 goroutine 完成
	// --------------------------
	if err := g.Wait(); err != nil {
		s.log.Error().Msgf("query user error: %v", err)
		return nil, exception.NewInternalServerError("%v", err)
	}

	return set, nil
}

// 查询用户列表
func (s *service) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.UserSet, error) {
	// // 0.判断这些要删除的用户是否存在
	queryReq := user.NewQueryUserDeleteRequest()
	queryReq.UserIds = req.UserIds
	s.log.Info().Msgf("this is user ids to delete: %s", req.UserIds)
	listusers, err := s.ListUser(ctx, queryReq)
	if err != nil {
		s.log.Error().Msgf("find user error, error is %s", err)
		return nil, exception.NewInternalServerError("find user error, error is %s", err)
	}
	usernames := []string{}
	if len(req.UserIds) != int(listusers.Total) {
		for _, user := range listusers.Items {
			usernames = append(usernames, user.Id)
		}
	} else {
		usernames = req.UserIds
	}
	//1.delete users from the user_ids
	if err := s.delete(ctx, listusers); err != nil {
		return nil, err
	}
	s.log.Info().Msgf("delete user %d, users list detail: %v", len(usernames), usernames)
	return listusers, nil
}
