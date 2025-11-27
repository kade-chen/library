package impl

import (
	"context"
	"math"

	"cloud.google.com/go/bigquery"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/iterator"
)

// 查询全部service 列表
func (s *service) QueryByDateProjectServicesAll(ctx context.Context, config *model.ProjectDataServiceSkuRequest) ([]model.ServicesList, error) {
	// 构造查询
	sql := s.queryByDateProjectSUSQL()
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
	q.Parameters = params

	// 执行查询
	it, err := q.Read(ctx)
	if err != nil {
		s.log.Error().Msgf("Failed to execute query, ERROR: %v", err)
		return nil, exception.NewInternalServerError("Failed to execute query, ERROR: %v", err)
	}

	var results []model.ServicesList
	// 遍历结果
	for {
		// var r Result
		var row model.ServicesList
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

// 1.默认查询所有的service and skus for project
func (s *service) QueryByDateService(ctx context.Context, config *model.ServiceDataRequest) ([]model.ServiceDateCost, error) {
	// startDate := "2025-10-01"
	// endDate := "2025-10-02"
	// projectIDs := []string{} // 空数组表示查询全部项目
	// projectIDs := []string{"yz-bx3-prod", "kade-poc", "zzshushu"} // 指定项目

	// 构造查询
	sql := s.queryByDateServiceSQL()
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

	if len(config.LabelKeys) > 0 {
		// 指定LabelKeys
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: config.LabelKeys})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "keys", Value: []string{}})
	}

	if len(config.LabelValues) > 0 {
		// 指定LabelValues
		params = append(params, bigquery.QueryParameter{Name: "value", Value: config.LabelValues})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "value", Value: []string{}})
	}

	if len(config.Region) > 0 {
		// 指定LabelValues
		params = append(params, bigquery.QueryParameter{Name: "region", Value: config.Region})
	} else {
		//查询全部
		params = append(params, bigquery.QueryParameter{Name: "region", Value: []string{}})
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

	var results []model.ServiceDateCost
	for {
		var row model.ServiceDateCost
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
