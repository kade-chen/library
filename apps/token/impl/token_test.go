package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	// _ "github.com/kade-chen/mcenter/apps" //registry all impl/api
	_ "github.com/kade-chen/google-billing-console/apps/auth/impl"
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
	req.Username = "kade@wondercloud.com"
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
	req := token.NewRevolkTokenRequest("KZd8hoMYnFDBP25HXswDJKwJ", "")
	req.ACCESS_TOKEN_NAME = "ACCESS_TOKEN_COOKIE_KEY"
	tk, row, err := impl.RevolkToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk, row)
	// t.Log(tk.Json())
}

func Test_Validate_Token(t *testing.T) {
	req := token.NewValidateTokenRequest()
	_, err := impl.ValicateToken(ctx, req)
	if err != nil {
		t.Fatal(err.Error())
	}
	// t.Log(err)
	// t.Log(tk.Json())
}

func Test_Update_Token(t *testing.T) {
	_, err := impl.UpdateToken(ctx, &token.UpdateTokenRequest{
		AccessToken:  "123",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJwbGF0Zm9ybSI6MCwic2NvcGUiOlsicGxhdGZvcm0uYWRtaW4iXSwiaXNzIjoid29uZGVyY2xvdWQuY29tIiwic3ViIjoia2FkZUB3b25kZXJjbG91ZC5jb20iLCJhdWQiOlsiZGV2LmJpbGxpbmcud29uZGVyY2xvdWQuY29tIl0sImV4cCI6MTc2NTk2NzA0NywiaWF0IjoxNzY1OTY3MDA3LCJqdGkiOiJqQkJjMWRhUHlCSDhpRnZsMmtCcGdRIn0.QyzRUcaOxX_rJ06xRrzl1P94wa6hqJVn6weiI_9t2X9z9fxudAFVaGQb__z7jV8vq_4c5xUiC8gg6UIhSo6sZLJsLKb3itZ44SGavGOmOBTTdwQT7pXN_sWx-CXBhzKUahMPOa_23ZoEo54FLvH3ccbFG3YmeOaBN6bclJCE12H9hCnFdZ5SLlTCaLBh5hdWzMGaynl7xy7UNj3e2KhTOQMI0MNC-9yIeaxU25GBkB3FNhTx4kk30EvwLFao3VEYqvrQ6-dJ24TKTT1fEjMXzB5y9e8e91eifsQKJoa3hw2sUPiwZKx5HHoNXZ3YDeE2OfHbPtbZ3n1TKB9NmJr-Qw",
	})
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
