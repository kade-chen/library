package impl

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/google-billing-console/apps/tools"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

// 1.默认查询所有的service and skus for project
func (s *service) QueryByDateProjectAll(ctx context.Context, config *project.ProjectDataConfig) ([]project.ProjectCost, error) {
	// startDate := "2025-10-01"
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByDateProjectSQL()
	q := s.bq.Query(sql)
	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
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

	var results []project.ProjectCost
	for {
		var row project.ProjectCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query completed")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		if !row.ProjectID.Valid {
			row.ProjectID.StringVal = "[Charges not specific to a project]"
			row.ProjectID.Valid = true
		}
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
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

// 2.自定义查询
// 2.1 全部service，指定sku
func (s *service) QueryByDateProjectServicesCustomSku(ctx context.Context, config *project.ProjectDataConfig) ([]project.ProjectCost, error) {
	// startDate := "2025-10-01"
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByDateProjectServicesCustomSkusSQL()
	q := s.bq.Query(sql)

	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
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

	var results []project.ProjectCost
	for {
		var row project.ProjectCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query finished")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		if !row.ProjectID.Valid {
			row.ProjectID.StringVal = "[Charges not specific to a project]"
			row.ProjectID.Valid = true
		}
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
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

// 2.自定义查询
// 2.2 指定service，全sku
func (s *service) QueryByDateProjectCustomServicesAllSkus(ctx context.Context, config *project.ProjectDataConfig) ([]project.ProjectCost, error) {
	// startDate := "2025-10-01"
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByDateProjectCustomServicesAllSkusSQL()
	q := s.bq.Query(sql)

	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
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
		return nil, exception.NewInternalServerError("Failed to execute query")
	}

	var results []project.ProjectCost
	for {
		var row project.ProjectCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query completed successfully")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		if !row.ProjectID.Valid {
			row.ProjectID.StringVal = "[Charges not specific to a project]"
			row.ProjectID.Valid = true
		}
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
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

// 2.自定义查询
// 2.2 指定service，指定sku
func (s *service) QueryByDateProjectCustomServicesCustomSkus(ctx context.Context, config *project.ProjectDataConfig) ([]project.ProjectCost, error) {
	// startDate := "2025-10-01"
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByDateProjectCustomServicesCustomSkusSQL()
	q := s.bq.Query(sql)

	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
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

	var results []project.ProjectCost
	for {
		var row project.ProjectCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msg("Query completed")
			break
		}
		if err != nil {
			return nil, exception.NewInternalServerError("Failed to iterate over query results: %v", err)
		}
		if !row.ProjectID.Valid {
			row.ProjectID.StringVal = "[Charges not specific to a project]"
			row.ProjectID.Valid = true
		}
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
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

func (s *service) QueryByDateProjectAllServicesAllSkus(ctx context.Context, config *model.ProjectDataRequest) (model.ByDateProjectAllServicesSkusList, error) {
	a, err := s.svcs.QueryByDateProjectServicesAll(ctx, config)
	if err != nil {
		return model.ByDateProjectAllServicesSkusList{}, err
	}
	b, err := s.skus.QueryByDateProjectSKUsAll(ctx, config)
	if err != nil {
		return model.ByDateProjectAllServicesSkusList{}, err
	}
	return model.ByDateProjectAllServicesSkusList{
		Services: a,
		Skus:     b,
	}, nil
}
