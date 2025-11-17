package proviatetoken_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/kade-chen/library/ioc"
// 	_ "github.com/kade-chen/mcenter/apps"
// 	"github.com/kade-chen/mcenter/apps/token"
// 	"github.com/kade-chen/mcenter/apps/token/provider"
// )

// var (
// 	impl1 provider.TokenIssuer
// 	ctx   = context.Background()
// )

// // 内部颁发测试
// func TestIssueToken(t *testing.T) {
// 	req := token.NewPrivateIssueTokenRequest("wg6H6Gn8vr2fc1sCDdQdrG1o", "test")
// 	req.Username = "kade.chen"

// 	// req.Password =
// 	t1, err := impl1.IssueToken(ctx, req)
// 	t.Log(t1, err)
// }

// func init() {
// 	req := ioc.NewLoadConfigRequest()
// 	req.ConfigFile.Enabled = true
// 	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/mcenter/etc/config.toml"
// 	ioc.DevelopmentSetup(req)
// 	impl1 = provider.GetTokenIssuer(token.GRANT_TYPE_PRIVATE_TOKEN)
// }
