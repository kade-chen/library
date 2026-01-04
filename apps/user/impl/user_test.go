package impl_test

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/user"
	_ "github.com/kade-chen/google-billing-console/apps/user/impl"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ctx  = context.Background()
	impl user.Service
)

// create user
func TestCreateUser(t *testing.T) {
	u, err := impl.CreateUser(ctx, &user.CreateUserRequest{
		Username: "kade1@vandercloud.com",
		Password: "123456",
		Domain:   []string{"vandercloud.com", "test.com", "test3.com"},
		// Type:     user.TYPE_SUPPER,
		Labels: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"cc": structpb.NewStringValue("value1"),
				"bb": structpb.NewNumberValue(123),
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(format.ToJSON(u))
}
func TestDescribeUser(t *testing.T) {
	req := user.NewDescriptUserRequestByName()
	// req.DescribeBy = user.DESCRIBE_BY_USER_ID
	// req.Id = "kadeqq11111"

	req.DescribeBy = user.DESCRIBE_BY_USER_NAME
	req.Domain = "wondercloud.com"
	req.Username = "kadeqq11111"
	a, err := impl.DescribeUser(ctx, req)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(a.Spec.Labels.Fields["key1"].GetStringValue())
	fmt.Println(format.ToJSON(a))
	// t.Log(a.Tojson1())
}

// list user
func TestListUser(t *testing.T) {
	req := user.NewQueryUserRequest(nil)
	// var myType user.TYPE = user.TYPE_SUB
	// req.Type = &myType
	// req.Domain = []string{"vandercloud.com","wondercloud3.com"}
	// req.UserIds = []string{"top@wondercloud.com","kade@wondercloud.com"}
	// req.Keywords = "kade"
	// req.Labels = &structpb.Struct{
	// 	Fields: map[string]*structpb.Value{
	// 		"cc": structpb.NewStringValue("value1"),
	// 		"key2": structpb.NewNumberValue(123),
	// 	},
	// }
	//总数
	req.Page.PageSize = 10
	//跳过多少个数据
	req.Page.Offset = 1
	a, err := impl.ListUser(ctx, req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(format.ToJSON(a))
	// t.Log(a.Tojson1())
}

// delete user
func TestDeleteUser(t *testing.T) {
	req := user.NewDeleteUserRequest()
	req.UserIds = []string{"top1@wondercloud.com", "top3@wondercloud.com", "top2@wondercloud.com", "kade"}
	a, err := impl.DeleteUser(ctx, req)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(a)
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(user.AppName).(user.Service)
}
