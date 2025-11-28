package impl

import (
	"context"
	"math"

	"cloud.google.com/go/bigquery"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/google-billing-console/tools/trances"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

// 查询全部service 列表
func (s *service) QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataServiceSkuRequest) ([]model.ServicesList, error) {
	if ctx.Value("trances_id") == nil {
		ctx = context.WithValue(context.Background(), "trances_id", trances.NewTraceID())
	}
	s.log.Info().Msgf("trances_id=%s, The User begins Query for UsageDateByDatesServiceSkusAPI", ctx.Value("trances_id"))
	// 构造查询
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL......", ctx.Value("trances_id"))
	sql := s.queryByDateProjectSUSQL()
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL successful", ctx.Value("trances_id"))

	q := s.bq.Query(sql)
	s.log.Info().Msgf("trances_id=%s, Configuring query parameters......", ctx.Value("trances_id"))

	partitionStartTime, partitionEndTime := tools.PartitionTime(config.StartDate, config.EndDate)
	s.log.Info().Msgf("trances_id=%s, start_date=%v, end_date=%v, PartitionStartTime=%v, PartitionEndTime=%v", ctx.Value("trances_id"), config.StartDate, config.EndDate, partitionStartTime, partitionEndTime)

	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
	}
	s.log.Info().Msgf("trances_id=%s, Retrieving project_ids......", ctx.Value("trances_id"))
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, project_ids=%v", ctx.Value("trances_id"), config.ProjectIDs)
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, project_ids is empty, query all project_ids", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}
	q.Parameters = params

	// 执行查询
	s.log.Info().Msgf("trances_id=%s, Executing query", ctx.Value("trances_id"))
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
		return nil, exception.NewInternalServerError("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
	}

	var results []model.ServicesList
	// 遍历结果
	for {
		// var r Result
		var row model.ServicesList
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msgf("trances_id=%s, Query completed for UsageDateByDatesServiceSkusAPI", ctx.Value("trances_id"))
			break
		}
		if err != nil {
			s.log.Error().Msgf("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
			return nil, exception.NewInternalServerError("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
		}
		results = append(results, row)
	}

	return results, err
}

// 1.默认查询所有的service and skus for project
func (s *service) QueryByDateService(ctx context.Context, config *model.ServiceDataRequest) ([]model.ServiceDateCost, error) {
	// 构造查询
	if ctx.Value("trances_id") == nil {
		ctx = context.WithValue(context.Background(), "trances_id", trances.NewTraceID())
	}
	s.log.Info().Msgf("trances_id=%s, The User begins Query for InvoiceMonthByDateServiceAPI", ctx.Value("trances_id"))
	// 构造查询
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL......", ctx.Value("trances_id"))

	// 构造查询
	sql := s.queryByDateServiceSQL()
	s.log.Info().Msgf("trances_id=%s, Configuring query parameters......", ctx.Value("trances_id"))

	q := s.bq.Query(sql)
	s.log.Info().Msgf("trances_id=%s, Configuring query parameters......", ctx.Value("trances_id"))

	partitionStartTime, partitionEndTime := tools.InvoiceMonthPartitionTime(config.StartDate, config.EndDate)
	s.log.Info().Msgf("trances_id=%s, start_date=%v, end_date=%v, PartitionStartTime=%v, PartitionEndTime=%v", ctx.Value("trances_id"), config.StartDate, config.EndDate, partitionStartTime, partitionEndTime)

	// 绑定参数
	params := []bigquery.QueryParameter{
		{Name: "start_date", Value: config.StartDate},
		{Name: "end_date", Value: config.EndDate},
		{Name: "PartitionStartTime", Value: partitionStartTime},
		{Name: "PartitionEndTime", Value: partitionEndTime},
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving project_ids......", ctx.Value("trances_id"))
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, project_ids=%v", ctx.Value("trances_id"), config.ProjectIDs)
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, project_ids is empty, query all project_ids", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving services_ids......", ctx.Value("trances_id"))
	if len(config.ServiceIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, services_ids=%v", ctx.Value("trances_id"), config.ServiceIDs)
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: config.ServiceIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, services_ids is empty, query all services_ids", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving skus_ids......", ctx.Value("trances_id"))
	if len(config.SkusIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, skus_ids=%v", ctx.Value("trances_id"), config.SkusIDs)
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: config.SkusIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, skus_ids is empty, query all skus_ids", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving label_keys......", ctx.Value("trances_id"))
	if len(config.LabelKeys) > 0 {
		// 指定LabelKeys
		s.log.Info().Msgf("trances_id=%s, label_keys=%v", ctx.Value("trances_id"), config.LabelKeys)
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: config.LabelKeys})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, label_keys is empty, query all label_keys", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving label_values......", ctx.Value("trances_id"))
	if len(config.LabelValues) > 0 {
		// 指定LabelValues
		s.log.Info().Msgf("trances_id=%s, label_values=%v", ctx.Value("trances_id"), config.LabelValues)
		params = append(params, bigquery.QueryParameter{Name: "value", Value: config.LabelValues})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, label_values is empty, query all label_values", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "value", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving region......", ctx.Value("trances_id"))
	if len(config.Region) > 0 {
		// 指定LabelValues
		s.log.Info().Msgf("trances_id=%s, region=%v", ctx.Value("trances_id"), config.Region)
		params = append(params, bigquery.QueryParameter{Name: "region", Value: config.Region})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, region is empty, query all region", ctx.Value("trances_id"))
		params = append(params, bigquery.QueryParameter{Name: "region", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving NegotiatedSavingsEnabled......", ctx.Value("trances_id"))
	s.log.Info().Msgf("trances_id=%s, NegotiatedSavingsEnabled=%v", ctx.Value("trances_id"), config.NegotiatedSavingsEnabled)

	var SavingsProgramsList []string

	s.log.Info().Msgf("trances_id=%s, Retrieving SavingsProgramsList......", ctx.Value("trances_id"))
	s.log.Info().Msgf("trances_id=%s, SavingsProgramsCommittedUsageDiscountEnabled=%v", ctx.Value("trances_id"), config.SavingsProgramsCommittedUsageDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, SavingsProgramsCommittedUsageDiscountDollarBaseEnabled=%v", ctx.Value("trances_id"), config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled)
	if config.SavingsProgramsCommittedUsageDiscountEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT")
	}
	if config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT_DOLLAR_BASE")
	}

	params = append(params, bigquery.QueryParameter{Name: "savings_programs_types", Value: SavingsProgramsList})

	s.log.Info().Msgf("trances_id=%s, Retrieving savings_programs_list......", ctx.Value("trances_id"))
	s.log.Info().Msgf("trances_id=%s, OtherSavingsFreeTierEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsFreeTierEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsPromotionEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsPromotionEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsSustainedUsageDiscountEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsSustainedUsageDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsResellerMarginEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsResellerMarginEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsDiscountEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsSubscriptionBenefitEnabled=%v", ctx.Value("trances_id"), config.OtherSavingsSubscriptionBenefitEnabled)
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
	s.log.Info().Msgf("trances_id=%s, Configuring acquisition complete", ctx.Value("trances_id"))
	q.Parameters = params

	// 执行查询
	s.log.Info().Msgf("trances_id=%s, Executing query", ctx.Value("trances_id"))
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
		return nil, exception.NewInternalServerError("trances_id=%s, Failed to execute query: %v", ctx.Value("trances_id"), err)
	}

	var results []model.ServiceDateCost
	for {
		var row model.ServiceDateCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msgf("trances_id=%s, Query completed for InvoiceMonthByDateServiceAPI", ctx.Value("trances_id"))
			break
		}
		if err != nil {
			s.log.Error().Msgf("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
			return nil, exception.NewInternalServerError("trances_id=%s, Failed to iterate over query results: %v", ctx.Value("trances_id"), err)
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
			s.log.Info().Msgf("trances_id=%s, TwoDecimalEnabled=%v", ctx.Value("trances_id"), config.TwoDecimalEnabled)
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
