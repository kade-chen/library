package user

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUserRequest
func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{
		Domain:          "kade-domain",
		Labels:          map[string]string{},
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
		SkipItems:    false,
		Labels:       map[string]string{},
		UserIds:      []string{},
		ExtraUserIds: []string{},
	}
}

// new QueryUserRequest request
func NewQueryUserRequest(req *QueryUserRequest) *QueryUserRequest {
	if req == nil {
		return &QueryUserRequest{
			// Page:         NewPageRequest(20, 1),
			Page:         &PageRequest{},
			SkipItems:    false,
			Labels:       map[string]string{},
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
		Page:         &PageRequest{},
		SkipItems:    false,
		Labels:       map[string]string{},
		UserIds:      []string{},
		ExtraUserIds: []string{},
	}
}

func (r *QueryUserRequest) FindFilter() bson.M {
	filter := bson.M{}

	if r.Domain != "" {
		filter["domain"] = r.Domain
	}
	if r.Provider != nil {
		filter["provider"] = r.Provider
	}
	if r.Type != nil {
		filter["type"] = r.Type
	}
	if len(r.UserIds) > 0 {
		filter["_id"] = bson.M{"$in": r.UserIds}
	}
	if r.Keywords != "" {
		filter["username"] = bson.M{"$regex": r.Keywords, "$options": "im"}
	}
	if r.Labels != nil {
		for k, v := range r.Labels {
			filter["labels."+k] = v
		}
	}

	if len(r.ExtraUserIds) > 0 {
		filter = bson.M{"$or": bson.A{
			filter,
			bson.M{"_id": bson.M{"$in": r.ExtraUserIds}},
		}}
	}

	return filter
}

// QueryUserRequest Options
func (r *QueryUserRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	// fmt.Println("--------", skip, pageSize)
	return &options.FindOptions{
		Sort: bson.D{
			{Key: "create_at", Value: -1},
		},
		Limit: &pageSize,
		Skip:  &skip,
	}
}

// NewDescriptUserRequestByName 查询详情请求
func NewDescriptUserRequestByName(username string) *DescribeUserRequest {
	return &DescribeUserRequest{
		Username:   username,
		DescribeBy: DESCRIBE_BY_USER_NAME,
	}
}

// NewDeleteUserRequest
func NewDeleteUserRequest() *DeleteUserRequest {
	return &DeleteUserRequest{
		UserIds: []string{},
	}
}
