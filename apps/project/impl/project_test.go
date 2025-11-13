package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx1  = context.Background()
	impl1 project.Service
)

func TestQueryByProjectAll(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataConfig
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01"} // 指定项目
	config.NegotiatedSavingsEnabled = true
	config.SavingsProgramsCommittedUsageDiscountEnabled = true
	config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled = true
	config.OtherSavingsFreeTierEnabled = true
	config.OtherSavingsPromotionEnabled = true
	config.OtherSavingsSustainedUsageDiscountEnabled = true
	config.OtherSavingsResellerMarginEnabled = true
	config.OtherSavingsDiscountEnabled = true
	config.OtherSavingsSubscriptionBenefitEnabled = true

	a, err := impl1.QueryByProjectAll(ctx1, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByProjectServicesCustomSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataConfig
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01"} // 指定项目
	config.NegotiatedSavingsEnabled = true
	config.SavingsProgramsCommittedUsageDiscountEnabled = true
	config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled = true
	config.OtherSavingsFreeTierEnabled = true
	config.OtherSavingsPromotionEnabled = true
	config.OtherSavingsSustainedUsageDiscountEnabled = true
	config.OtherSavingsResellerMarginEnabled = true
	config.OtherSavingsDiscountEnabled = true
	config.OtherSavingsSubscriptionBenefitEnabled = true
	config.SkusIDs = []string{"984A-1F27-2D1F"}
	a, err := impl1.QueryByProjectServicesCustomSku(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByProjectCustomServicesAllSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataConfig
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01"} // 指定项目
	config.NegotiatedSavingsEnabled = true
	config.SavingsProgramsCommittedUsageDiscountEnabled = true
	config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled = true
	config.OtherSavingsFreeTierEnabled = true
	config.OtherSavingsPromotionEnabled = true
	config.OtherSavingsSustainedUsageDiscountEnabled = true
	config.OtherSavingsResellerMarginEnabled = true
	config.OtherSavingsDiscountEnabled = true
	config.OtherSavingsSubscriptionBenefitEnabled = true
	config.ServiceIDs = []string{"6F81-5844-456A"}
	a, err := impl1.QueryByProjectCustomServicesAllSkus(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByProjectCustomServicesCustomSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataConfig
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01"} // 指定项目
	config.NegotiatedSavingsEnabled = true
	config.SavingsProgramsCommittedUsageDiscountEnabled = true
	config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled = true
	config.OtherSavingsFreeTierEnabled = true
	config.OtherSavingsPromotionEnabled = true
	config.OtherSavingsSustainedUsageDiscountEnabled = true
	config.OtherSavingsResellerMarginEnabled = true
	config.OtherSavingsDiscountEnabled = true
	config.OtherSavingsSubscriptionBenefitEnabled = true
	config.ServiceIDs = []string{"6F81-5844-456A"}
	config.SkusIDs = []string{"6CB7-B05F-97AD", "DE9E-AFBC-A15A"}
	a, err := impl1.QueryByProjectCustomServicesCustomSkus(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}


func Test111(t *testing.T) {
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl1 = ioc.Controller().Get(project.AppName).(project.Service)
}
