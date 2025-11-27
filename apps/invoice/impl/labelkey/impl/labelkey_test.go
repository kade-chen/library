package impl_test

import (
	"context"
	"fmt"
	"testing"

	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/labelkey"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl invoice.LabelKeyService
)

func TestSKU(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	a, err := impl.QueryByInvoiceMonthProjectLabelKeyAll(ctx, &model.InvoiceMonthProjectLabelKeyRequest{
		StartDate:  "202510",
		EndDate:    "202510",
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
	impl = ioc.Controller().Get(labelkey.AppName).(invoice.LabelKeyService)
}
