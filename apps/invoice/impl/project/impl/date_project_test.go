package impl_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/project"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps"
)

var (
	ctx  = context.Background()
	impl invoice.ProjectService
)

func TestQueryByDateProject1(t *testing.T) {
	s, e, err := PartitionTime("202508", "202510")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s, e)
}

func PartitionTime(startStr, endStr string) (string, string, error) {
	// 解析 YYYYMM → time.Time（默认取每月1号）
	start, err := time.Parse("200601", startStr)
	if err != nil {
		return "", "", fmt.Errorf("start_date invalid: %w", err)
	}

	end, err := time.Parse("200601", endStr)
	if err != nil {
		return "", "", fmt.Errorf("end_date invalid: %w", err)
	}

	// prev_start = start - 1 个月（取每月 1 号）
	prevStart := start.AddDate(0, -1, 0)
	prevStartStr := fmt.Sprintf("%04d-%02d-01", prevStart.Year(), prevStart.Month())

	// prev_end = end - 1 个月的“月底”
	// 技巧：跳到下个月 1 号再往前减一天，就是这个月的月底
	prevEndMonth := end.AddDate(0, 0, 0)
	nextMonth := time.Date(prevEndMonth.Year(), prevEndMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	prevEnd := nextMonth.Add(-24 * time.Hour)
	prevEndStr := prevEnd.Format("2006-01-02")

	return prevStartStr, prevEndStr, nil
}

func TestQueryByDateProject(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataRequest
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
	config.TwoDecimalEnabled = false
	a, err := impl.QueryByDateProject(ctx, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(format.ToJSON(a))
}

func TestQueryByDateProjectAllServicesAllSkus(t *testing.T) {
	fmt.Println(ioc.Controller().List())
	var config model.ProjectDataServiceSkuRequest
	config.StartDate = "2025-10-01"
	config.EndDate = "2025-10-02"
	config.ProjectIDs = []string{"tools-orion", "chat-prod-404613", "sw-pro-01", "ffalcon-hw-01"} // 指定项目
	a, err := impl.QueryByDateProjectAllServicesAllSkus(ctx, &config)
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
	impl = ioc.Controller().Get(project.AppName).(invoice.ProjectService)
}
