package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/services"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl services.Service
)

func TestMain1(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	a, _ := impl.QueryByDateProjectServicesAll(ctx, &model.ProjectDataRequest{})
	fmt.Println(format.ToJSON(a))
}

func TestByDateService1(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ServiceDataConfig
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01", "kade-poc"} // 指定项目
	config.NegotiatedSavingsEnabled = true
	config.SavingsProgramsCommittedUsageDiscountEnabled = true
	config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled = true
	config.OtherSavingsFreeTierEnabled = true
	config.OtherSavingsPromotionEnabled = true
	config.OtherSavingsSustainedUsageDiscountEnabled = true
	config.OtherSavingsResellerMarginEnabled = true
	config.OtherSavingsDiscountEnabled = true
	config.OtherSavingsSubscriptionBenefitEnabled = true
	//查询全部
	//查询servoce
	config.ServiceIDs = []string{"6F81-5844-456A"}
	//sku
	config.SkusIDs = []string{"DE9E-AFBC-A15A", "6CB7-B05F-97AD"}
	//查询service/sku
	// config.ServiceIDs = []string{"6F81-5844-456A"}
	// config.SkusIDs = []string{"DE9E-AFBC-A15A", "6CB7-B05F-97AD"}
	a, _ := impl.QueryByDateService(ctx, &config)
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(services.AppName).(services.Service)
}
