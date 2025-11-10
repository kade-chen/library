package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/google-billing-console/apps/services"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl services.Service
)

func TestMain(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	a, _ := impl.QueryByDateProjectServicesAll(ctx, &project.ProjectDataConfig{})
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(services.AppName).(services.Service)
}
