package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/billing"
	"github.com/kade-chen/library/ioc"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl billing.Service
)

func TestMain(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	impl.QueryByDateProject(ctx, "1")
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/mcenter/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(billing.AppName).(billing.Service)
}
