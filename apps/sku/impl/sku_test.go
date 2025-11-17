package impl_test

import (
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

func TestSku(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.SkuConfig
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
	a, _ := impl.QueryBySku(ctx, &config)
	fmt.Println(format.ToJSON(a))
}
