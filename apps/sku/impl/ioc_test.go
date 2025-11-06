package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/sku"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl sku.Service
)

func TestSKU(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	a, _ := impl.QueryByDateProjectSKUsAll(ctx, "1")
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/mcenter/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(sku.AppName).(sku.Service)
}
