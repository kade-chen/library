package impl

import (
	"context"
	"fmt"
	"math"

	"cloud.google.com/go/bigquery"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByDateProjectSKUsAll(ctx context.Context, config *model.ProjectDataServiceSkuRequest) ([]model.SkusList, error) {
	// config.StartDate = "2025-10-01"
	// config.EndDate = "2025-10-02"
	// // projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateProjectSKUsSQL()
	q := s.bq.Query(sql)

	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	s.log.Debug().Msgf("partitionStartTime:%s partitionEndTime:%s", partitionStartTime, partitionEndTime)
	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
	}
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}
	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		return nil, exception.NewInternalServerError("Failed to execute query, ERROR: %v", err)
	}

	var results []model.SkusList
	// 遍历结果
	for {
		// var r Result
		var row model.SkusList
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query finished")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		results = append(results, row)
	}

	return results, err
}

func (s *service) QueryByDateSku(ctx context.Context, config *model.SkuDataRequest) ([]model.SkuDateCost, error) {
	// config.StartDate = "2025-10-01"
	// config.EndDate = "2025-10-02"
	// // projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateSkuSQL()
	q := s.bq.Query(sql)
	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	fmt.Println("partitionStartTime:", partitionStartTime, "partitionEndTime:", partitionEndTime)
	fmt.Println("StartTime:", config.StartDate, "EndTime:", config.EndDate)
	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
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

	var results []model.SkuDateCost
	for {
		var row model.SkuDateCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query completed")
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
		results = append(results, row)
	}
	return results, nil
}

func (s *service) QueryByDateSkuHeru(ctx context.Context, config *model.SkuDataRequest) ([]model.AlibabaHehuSkuDateCost, error) {
	// config.StartDate = "2025-10-01"
	// config.EndDate = "2025-10-02"
	// // projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目
	// 构造查询
	sql := s.queryByDateSkuSQL1()
	q := s.bq.Query(sql)
	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	fmt.Println("partitionStartTime:", partitionStartTime, "partitionEndTime:", partitionEndTime)
	fmt.Println("StartTime:", config.StartDate, "EndTime:", config.EndDate)
	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
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

	var results []model.AlibabaHehuSkuDateCost
	for {
		var row model.AlibabaHehuSkuDateCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query completed")
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
		// if !config.SavingsProgramsCommittedUsageDiscountEnabled  && !config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled  {
		// 	row.SavingsPrograms.Float64 = 0.00
		// 	row.NegotiatedSavings.Valid = true
		// }
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
