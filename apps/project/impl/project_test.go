package impl_test

import (
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
	"github.com/kade-chen/google-billing-console/apps/common/model"
)

func TestQueryByProjectCustomServicesCustomSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectConfig
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
	config.ServiceIDs = []string{"6F81-5844-456A"}
	config.SkusIDs = []string{"6CB7-B05F-97AD", "DE9E-AFBC-A15A"}
	a, err := impl.QueryByProject(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}
