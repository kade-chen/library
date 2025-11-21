package impl

import (
	"context"
	"math"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByService(ctx context.Context, config *model.ServiceConfig) ([]model.ServiceCost, error) {
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByServiceSQL()
	q := s.bq.Query(sql)

	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
	prev_start, prev_end := tools.PartitionPrevDates(config.StartDate, config.EndDate)

	prev_PartitionStartTime, prev_PartitionEndTime := tools.PartitionTime(prev_start, prev_end)
	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "prev_start", Value: prev_start},
		{Name: "prev_end", Value: prev_end},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
		{Name: "prev_PartitionStartTime", Value: prev_PartitionStartTime},
		{Name: "prev_PartitionEndTime", Value: prev_PartitionEndTime},
	}
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}

	if len(config.ServiceIDs) > 0 {
		// 指定projectt
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: config.ServiceIDs})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: []string{}})
	}

	if len(config.SkusIDs) > 0 {
		// 指定projectt
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: config.SkusIDs})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: []string{}})
	}
	var SavingsProgramsList []string

	if config.SavingsProgramsCommittedUsageDiscountEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT")
	}
	if config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT_DOLLAR_BASE")
	}

	params = append(params, bigquery.QueryParameter{Name: "savings_programs_types", Value: SavingsProgramsList})

	var OtherSavingsList []string
	//https://cloud.google.com/billing/docs/how-to/cost-table#columns_in_the_cost_table
	if config.OtherSavingsFreeTierEnabled {
		OtherSavingsList = append(OtherSavingsList, "FREE_TIER")
	}
	if config.OtherSavingsPromotionEnabled {
		OtherSavingsList = append(OtherSavingsList, "PROMOTION")
	}
	if config.OtherSavingsSustainedUsageDiscountEnabled {
		OtherSavingsList = append(OtherSavingsList, "SUSTAINED_USAGE_DISCOUNT")
	}
	if config.OtherSavingsResellerMarginEnabled {
		OtherSavingsList = append(OtherSavingsList, "RESELLER_MARGIN")
	}
	if config.OtherSavingsDiscountEnabled {
		OtherSavingsList = append(OtherSavingsList, "DISCOUNT")
	}
	if config.OtherSavingsSubscriptionBenefitEnabled {
		OtherSavingsList = append(OtherSavingsList, "SUBSCRIPTION_BENEFIT")
	}

	params = append(params, bigquery.QueryParameter{Name: "other_savings_types", Value: OtherSavingsList})
	// 	if (!config.SavingsProgramsCommittedUsageDiscountEnabled && !config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled) && (!config.SavingsProgramsCommittedUsageDiscountEnabled || !config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled) {
	// 		params = append(params, bigquery.QueryParameter{Name: "savings_programs_types", Value: []string{"COMMITTED_USAGE_DISCOUNT", "COMMITTED_USAGE_DISCOUNT_DOLLAR_BASE"}})
	// 	} else {
	// 		params = append(params, bigquery.QueryParameter{Name: "savings_programs_types", Value: []string{}})
	// 	}
	// }

	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("Failed to execute query: %v", err)
		return nil, exception.NewInternalServerError("Failed to execute query: %v", err)
	}

	var results []model.ServiceCost
	for {
		var row model.ServiceCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query finished")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		// if !row.ProjectID.Valid {
		// 	row.ProjectID.StringVal = "[Charges not specific to a project]"
		// 	row.ProjectID.Valid = true
		// }
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
		}
		if config.TwoDecimalEnabled {
			row.UsageCost.Float64 = math.Round(row.UsageCost.Float64*100) / 100
			row.UsageCost.Valid = true
			row.NegotiatedSavings.Float64 = math.Round(row.NegotiatedSavings.Float64*100) / 100
			row.UsageCost.Valid = true
			row.SavingsPrograms.Float64 = math.Round(row.SavingsPrograms.Float64*100) / 100
			row.UsageCost.Valid = true
			row.OtherSavings.Float64 = math.Round(row.OtherSavings.Float64*100) / 100
			row.UsageCost.Valid = true
			row.SubTotal.Float64 = math.Round(row.SubTotal.Float64*100) / 100
			row.UsageCost.Valid = true
		}
		// fmt.Printf("%s | %s | %s | %s | %s | %s\n",
		// 	format.FormatAny(row["usage_date"]),
		// 	format.FormatAny(row["project_id"]),
		// 	format.FormatAny(row["invoice_cost"]),
		// 	format.FormatAny(row["invoice_cost_at_list_abs"]),
		// 	format.FormatAny(row["cost_at_list"]),
		// 	format.FormatAny(row["Usage_Cost"]),
		// )
		results = append(results, row)
	}
	return results, nil
}
