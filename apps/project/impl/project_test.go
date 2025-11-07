package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl project.Service
)

func TestQueryByDateProjectAll(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config project.ProjectDataConfig
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

	a, err := impl.QueryByDateProjectAll(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByDateProjectServicesCustomSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config project.ProjectDataConfig
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
	config.SkusIDs = []string{"4111-7FF1-D50A"}
	a, err := impl.QueryByDateProjectServicesCustomSku(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByDateProjectCustomServicesAllSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config project.ProjectDataConfig
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
	a, err := impl.QueryByDateProjectCustomServicesAllSkus(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByDateProjectCustomServicesCustomSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config project.ProjectDataConfig
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
	a, err := impl.QueryByDateProjectCustomServicesCustomSkus(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(project.AppName).(project.Service)
}
