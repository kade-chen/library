package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	// _ "github.com/kade-chen/mcenter/apps" //registry all impl/api
	_ "github.com/kade-chen/google-billing-console/apps/domain/impl"
	"github.com/kade-chen/google-billing-console/apps/token"
	_ "github.com/kade-chen/google-billing-console/apps/token/impl"
	_ "github.com/kade-chen/google-billing-console/apps/token/provider/all"
	_ "github.com/kade-chen/google-billing-console/apps/user/impl"
)

var (
	ctx  = context.Background()
	impl token.Service
)

func Test_Issue_Token_PassWord(t *testing.T) {
	req := token.NewIssueTokenRequest()
	req.Username = "kade"
	req.Password = "123456"
	// req.GrantType = token.GRANT_TYPE_PRIVATE_TOKEN
	tk, err := impl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tk.Json())
}

func Test_Issue_Token(t *testing.T) {
	req := token.NewIssueTokenRequest()
	req.Username = "admin"
	req.Password = "123456"
	req.GrantType = token.GRANT_TYPE_PRIVATE_TOKEN
	tk, err := impl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
	// t.Log(tk.Json())
}

func Test_Revoke_Token(t *testing.T) {
	req := token.NewRevolkTokenRequest("gm8buLsAZoXc5pkRR34QQoj6", "")
	req.ACCESS_TOKEN_NAME = "ACCESS_TOKEN_COOKIE_KEY"
	tk, err := impl.RevolkToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
	// t.Log(tk.Json())
}

func Test_Validate_Token(t *testing.T) {
	req := token.NewValidateTokenRequest("BUxe2dhJpzk8PyRsapF6AA5U")
	_, err := impl.ValicateToken(ctx, req)
	if err != nil {
		t.Fatal(err.Error())
	}
	// t.Log(err)
	// t.Log(tk.Json())
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(token.AppName).(token.Service)
}
