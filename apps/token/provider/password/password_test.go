package password_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/kade-chen/google-billing-console/apps"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"
)

var (
	impl provider.TokenIssuer
	ctx  = context.Background()
)

func TestIssueToken(t *testing.T) {
	// req := token.NewPasswordIssueTokenRequest("admin", "123456")
	tk, err := impl.IssueToken(ctx, &token.IssueTokenRequest{Username: "kade@wondercloud.com", Password: "123456"})
	if err != nil {
		fmt.Println(format.ToJSON(err))
		t.Errorf("issue token failed: %v", err)
		return
	}
	fmt.Println(tk.Json())
}

func TestMain1(t *testing.T) {
	fmt.Println(time.Unix(1731983085, 0))
	expiredAt := time.Unix(1731983085, 0).Add(time.Duration(1) * time.Hour * 24)
	fmt.Println(expiredAt)
	fmt.Println(time.Unix(1731983085, 0).Sub(expiredAt).Hours())
}

func init() {
	// os.Setenv("DEBUG", "true") //this debug conflicts with another vs debug
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = provider.GetTokenIssuer(token.GRANT_TYPE_PASSWORD)
}
