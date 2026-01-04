package user

import (
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// CreateUserRequest
func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{
		// Domain: "kade-domain",
		// Labels:          map[string]string{},
		Labels: &structpb.Struct{
			// Fields: map[string]*structpb.Value{
			// 	"key1": structpb.NewStringValue("value1"),
			// 	"key2": structpb.NewNumberValue(123),
			// },
		},
		Feishu:          &Feishu{},
		Dingding:        &DingDing{},
		Wechatwork:      &WechatWork{},
		Profile:         &Profile{},
		UseFullNamedUid: true,
	}
}

// NewQueryUserRequest 列表查询请求
func NewQueryUserRequestq() *QueryUserRequest {
	return &QueryUserRequest{
		Page: &PageRequest{
			PageSize:   uint64(20),
			PageNumber: uint64(1),
		},
		SkipItems: false,
		Labels: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"key1": structpb.NewStringValue("value1"),
				"key2": structpb.NewNumberValue(123),
			},
		},
		UserIds:      []string{},
		ExtraUserIds: []string{},
	}
}

// new QueryUserRequest request
func NewQueryUserRequest(req *QueryUserRequest) *QueryUserRequest {
	if req == nil {
		return &QueryUserRequest{
			// Page:         NewPageRequest(20, 1),
			Page:      &PageRequest{},
			Domain:    []string{},
			SkipItems: false,
			Labels: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					// "key1": structpb.NewStringValue("value1"),
					// "key2": structpb.NewNumberValue(123),
				},
			},
			UserIds:      []string{},
			ExtraUserIds: []string{},
		}
	}
	return req
}

// NewQueryUserRequest 列表查询请求
func NewQueryUserDeleteRequest() *QueryUserRequest {
	return &QueryUserRequest{
		// Page:         NewPageRequest(20, 1),
		Page:      &PageRequest{},
		SkipItems: false,
		Labels:    &structpb.Struct{
			// Fields: map[string]*structpb.Value{
			// 	"key1": structpb.NewStringValue("value1"),
			// 	"key2": structpb.NewNumberValue(123),
			// },
		},
		UserIds:      []string{},
		ExtraUserIds: []string{},
	}
}

func (r *QueryUserRequest) PageSQL() (string, []bigquery.QueryParameter) {
	conditions := []string{}
	params := []bigquery.QueryParameter{}

	// LIMIT
	if r.Page.PageSize > 0 {
		conditions = append(conditions, " LIMIT @limit")
		params = append(params,
			bigquery.QueryParameter{Name: "limit", Value: int64(r.Page.PageSize)},
		)
	}

	// OFFSET
	if r.Page.Offset > 0 {
		conditions = append(conditions, " OFFSET @offset")
		params = append(params,
			bigquery.QueryParameter{Name: "offset", Value: int64(r.Page.Offset)},
		)
	}
	return strings.Join(conditions, " "), params
}

func (r *QueryUserRequest) WhereSQL() (string, []bigquery.QueryParameter) {
	conditions := []string{}
	params := []bigquery.QueryParameter{}

	// domain
	if len(r.Domain) > 0 {
		conditions = append(conditions, `
        EXISTS (
            SELECT 1
            FROM UNNEST(spec.domain) d
            WHERE d IN UNNEST(@domain)
        )
    `)
		params = append(params, bigquery.QueryParameter{Name: "domain", Value: r.Domain})
	}

	// provider
	if r.Provider != nil {
		conditions = append(conditions, "spec.provider = @provider")
		params = append(params, bigquery.QueryParameter{Name: "provider", Value: int64(*r.Provider)})
	}

	// type
	if r.Type != nil {
		conditions = append(conditions, "spec.type = @type")
		params = append(params, bigquery.QueryParameter{Name: "type", Value: int64(*r.Type)})
	}

	// userIds
	if len(r.UserIds) > 0 {
		conditions = append(conditions, "id IN UNNEST(@user_ids)")
		params = append(params, bigquery.QueryParameter{Name: "user_ids", Value: r.UserIds})
	}

	// username regex → REGEXP_CONTAINS
	if r.Keywords != "" {
		conditions = append(conditions, "REGEXP_CONTAINS(spec.username, @keywords)")
		params = append(params, bigquery.QueryParameter{Name: "keywords", Value: r.Keywords})
	}

	if len(r.Labels.Fields) > 0 {
		// labels.x = y
		// 只支持 JSON_EXTRACT_SCALAR，所以 v 必须是 string
		if r.Labels != nil {
			for k, v := range r.Labels.Fields {
				paramName := "label_" + k
				conditions = append(conditions,
					fmt.Sprintf("JSON_VALUE(spec.labels, '$.%s') = @%s", k, paramName),
				)
				// ⬅⬅ BigQuery 参数必须是基本类型
				params = append(params, bigquery.QueryParameter{
					Name:  paramName,
					Value: v.GetStringValue(),
				})
			}
		}
	}

	// ExtraUserIds OR logic
	if len(r.ExtraUserIds) > 0 {
		left := "TRUE"
		if len(conditions) > 0 {
			left = "(" + strings.Join(conditions, " AND ") + ")"
		}

		right := "`_id` IN UNNEST(@extra_user_ids)"
		params = append(params, bigquery.QueryParameter{Name: "extra_user_ids", Value: r.ExtraUserIds})

		conditions = []string{
			fmt.Sprintf("(%s OR %s)", left, right),
		}
	}

	if len(conditions) == 0 {
		return "", params
	}

	return "WHERE " + strings.Join(conditions, " AND "), params
}

// NewDescriptUserRequestByName 查询详情请求
func NewDescriptUserRequestByName() *DescribeUserRequest {
	return &DescribeUserRequest{
		// Id:         id,
		// Username:   username,
		// DescribeBy: DESCRIBE_BY_USER_NAME,
	}
}

// NewDeleteUserRequest
func NewDeleteUserRequest() *DeleteUserRequest {
	return &DeleteUserRequest{
		UserIds: []string{},
	}
}
