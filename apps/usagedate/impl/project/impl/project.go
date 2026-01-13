package impl

import (
	"context"
	"math"

	"cloud.google.com/go/bigquery"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

func (s *service) QueryByProject(ctx context.Context, config *model.ProjectRequest) ([]model.ProjectCost, error) {
	// if trancesID == nil {
	// 	ctx = context.WithValue(context.Background(), "trances_id", trances.NewTraceID())
	// }
	trancesID := ctx.Value("claims").(*authModel.TokenAuthMiddleware).TrancesID
	s.log.Info().Msgf("trances_id=%s, The User begins Query for UsageDateByPojectAPI", trancesID)
	// 构造查询
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL......", trancesID)
	sql := s.queryByProjectSQL1(config.OrganizationBqTable)
	s.log.Info().Msgf("trances_id=%s, Retrieving initialization SQL successful", trancesID)
	q := s.bq.Query(sql)
	s.log.Info().Msgf("trances_id=%s, Configuring query parameters......", trancesID)

	partitionStartTime, partitionEndTime := tools.CustomPartitionTime(config.StartDate, config.EndDate, 2)
	prev_start, prev_end := tools.PartitionPrevDates(config.StartDate, config.EndDate)
	s.log.Info().Msgf("trances_id=%s, start_date=%v, end_date=%v, PartitionStartTime=%v, PartitionEndTime=%v", trancesID, config.StartDate, config.EndDate, partitionStartTime, partitionEndTime)

	prev_PartitionStartTime, prev_PartitionEndTime := tools.CustomPartitionTime(prev_start, prev_end, 2)
	s.log.Info().Msgf("trances_id=%s, prev_start_date=%v, prev_end_date=%v, prev_PartitionStartTime=%v, prev_PartitionEndTime=%v", trancesID, prev_start, prev_end, prev_PartitionStartTime, prev_PartitionEndTime)
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

	s.log.Info().Msgf("trances_id=%s, Retrieving project_ids......", trancesID)
	if len(config.ProjectIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, project_ids=%v", trancesID, config.ProjectIDs)
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: config.ProjectIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, project_ids is empty, query all project_ids", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "project_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving services_ids......", trancesID)
	if len(config.ServiceIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, services_ids=%v", trancesID, config.ServiceIDs)
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: config.ServiceIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, services_ids is empty, query all services_ids", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "services_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving skus_ids......", trancesID)
	if len(config.SkusIDs) > 0 {
		// 指定projectt
		s.log.Info().Msgf("trances_id=%s, skus_ids=%v", trancesID, config.SkusIDs)
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: config.SkusIDs})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, skus_ids is empty, query all skus_ids", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "skus_ids", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving label_keys......", trancesID)
	if len(config.LabelKeys) > 0 {
		// 指定LabelKeys
		s.log.Info().Msgf("trances_id=%s, label_keys=%v", trancesID, config.LabelKeys)
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: config.LabelKeys})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, label_keys is empty, query all label_keys", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving label_values......", trancesID)
	if len(config.LabelValues) > 0 {
		// 指定LabelValues
		s.log.Info().Msgf("trances_id=%s, label_values=%v", trancesID, config.LabelValues)
		params = append(params, bigquery.QueryParameter{Name: "value", Value: config.LabelValues})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, label_values is empty, query all label_values", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "value", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving region......", trancesID)
	if len(config.Region) > 0 {
		// 指定LabelValues
		s.log.Info().Msgf("trances_id=%s, region=%v", trancesID, config.Region)
		params = append(params, bigquery.QueryParameter{Name: "region", Value: config.Region})
	} else {
		//查询全部
		s.log.Info().Msgf("trances_id=%s, region is empty, query all region", trancesID)
		params = append(params, bigquery.QueryParameter{Name: "region", Value: []string{}})
	}

	s.log.Info().Msgf("trances_id=%s, Retrieving NegotiatedSavingsEnabled......", trancesID)
	s.log.Info().Msgf("trances_id=%s, NegotiatedSavingsEnabled=%v", trancesID, config.NegotiatedSavingsEnabled)

	var SavingsProgramsList []string

	s.log.Info().Msgf("trances_id=%s, Retrieving SavingsProgramsList......", trancesID)
	s.log.Info().Msgf("trances_id=%s, SavingsProgramsCommittedUsageDiscountEnabled=%v", trancesID, config.SavingsProgramsCommittedUsageDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, SavingsProgramsCommittedUsageDiscountDollarBaseEnabled=%v", trancesID, config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled)
	if config.SavingsProgramsCommittedUsageDiscountEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT")
	}
	if config.SavingsProgramsCommittedUsageDiscountDollarBaseEnabled {
		SavingsProgramsList = append(SavingsProgramsList, "COMMITTED_USAGE_DISCOUNT_DOLLAR_BASE")
	}

	params = append(params, bigquery.QueryParameter{Name: "savings_programs_types", Value: SavingsProgramsList})

	s.log.Info().Msgf("trances_id=%s, Retrieving savings_programs_list......", trancesID)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsFreeTierEnabled=%v", trancesID, config.OtherSavingsFreeTierEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsPromotionEnabled=%v", trancesID, config.OtherSavingsPromotionEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsSustainedUsageDiscountEnabled=%v", trancesID, config.OtherSavingsSustainedUsageDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsResellerMarginEnabled=%v", trancesID, config.OtherSavingsResellerMarginEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsDiscountEnabled=%v", trancesID, config.OtherSavingsDiscountEnabled)
	s.log.Info().Msgf("trances_id=%s, OtherSavingsSubscriptionBenefitEnabled=%v", trancesID, config.OtherSavingsSubscriptionBenefitEnabled)
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
	s.log.Info().Msgf("trances_id=%s, Configuring acquisition complete", trancesID)
	q.Parameters = params

	// 执行查询
	s.log.Info().Msgf("trances_id=%s, Executing query", trancesID)
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("trances_id=%s, Failed to execute query: %v", trancesID, err)
		return nil, exception.NewInternalServerError("trances_id=%s, Failed to execute query: %v", trancesID, err)
	}

	var results []model.ProjectCost
	for {
		var row model.ProjectCost
		err := it.Next(&row)
		if err == iterator.Done {
			s.log.Info().Msgf("trances_id=%s, Query completed for UsageDateByPojectAPI", trancesID)
			break
		}
		if err != nil {
			s.log.Error().Msgf("trances_id=%s, Failed to iterate over query results: %v", trancesID, err)
			return nil, exception.NewInternalServerError("trances_id=%s, Failed to iterate over query results: %v", trancesID, err)
		}
		if !row.ProjectID.Valid {
			s.log.Info().Msgf("trances_id=%s, ProjectID is null, setting to '[Charges not specific to a project]", trancesID)
			row.ProjectID.StringVal = "[Charges not specific to a project]"
			row.ProjectID.Valid = true
		}
		if !config.NegotiatedSavingsEnabled {
			row.NegotiatedSavings.Float64 = 0.00
			row.NegotiatedSavings.Valid = true
		}
		if config.TwoDecimalEnabled {
			s.log.Info().Msgf("trances_id=%s, TwoDecimalEnabled=%v", trancesID, config.TwoDecimalEnabled)
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
