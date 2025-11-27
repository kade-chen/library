package impl_test

import (
	"context"
	"fmt"
	"testing"

	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
	"github.com/kade-chen/google-billing-console/apps/usagedate"
	"github.com/kade-chen/google-billing-console/apps/usagedate/impl/labelkey"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl usagedate.LabelKeyService
)

func TestSKU(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	a, err := impl.QueryByUsageDatProjectLabelKeyAll(ctx, &model.UsageDateProjectLabelKeyRequest{
		StartDate:  "2025-10-01",
		EndDate:    "2025-10-01",
		ProjectIDs: []string{"gen-lang-client-0334376452"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(labelkey.AppName).(usagedate.LabelKeyService)
}
