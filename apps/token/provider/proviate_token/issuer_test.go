package proviatetoken_test

import (
	"context"
	"testing"

	_ "github.com/kade-chen/google-billing-console/apps"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/library/ioc"
)

var (
	impl1 provider.TokenIssuer
	ctx   = context.Background()
)

// 内部颁发测试
func TestIssueToken(t *testing.T) {
	req := token.NewPrivateIssueTokenRequest("wg6H6Gn8vr2fc1sCDdQdrG1o", "test")
	req.Username = "kade.chen"

	// req.Password =
	t1, err := impl1.IssueToken(ctx, req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(t1, err)
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl1 = provider.GetTokenIssuer(token.GRANT_TYPE_PRIVATE_TOKEN)
}
