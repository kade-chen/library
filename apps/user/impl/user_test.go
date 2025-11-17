package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/user"
	_ "github.com/kade-chen/google-billing-console/apps/user/impl"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"
)

var (
	ctx  = context.Background()
	impl user.Service
)

// create user
func TestCreateUser(t *testing.T) {
	u, err := impl.CreateUser(ctx, &user.CreateUserRequest{
		Username: "kade",
		Password: "123456",
		Domain:   "wondercloud.com",
		Type:     user.TYPE_SUB,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(format.ToJSON(u))
}

// query user
func TestQueryUser(t *testing.T) {
	req := user.NewQueryUserRequest(nil)
	var myType user.TYPE = user.TYPE_SUB
	req.Type = &myType
	a, err := impl.ListUser(ctx, req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
	// t.Log(a.Tojson1())
}

// delete user
func TestDeleteUser(t *testing.T) {
	req := user.NewDeleteUserRequest()
	req.UserIds = []string{"oo@kade-domain"}
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
