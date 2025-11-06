package project

import (
	"context"

	"cloud.google.com/go/bigquery"
)

const (
	AppName = "bigquery"
)

type Service interface {
	// QueryByProject(ctx context.Context, query string) error
	QueryByDateProjectAll(context.Context, *ProjectDataConfig) ([]ProjectCost, error)
	QueryByDateProjectServicesCustomSku(context.Context, *ProjectDataConfig) ([]ProjectCost, error)
	QueryByDateProjectCustomServicesAllSkus(context.Context, *ProjectDataConfig) ([]ProjectCost, error)
	QueryByDateProjectCustomServicesCustomSkus(context.Context, *ProjectDataConfig) ([]ProjectCost, error)
	// QueryBySku(ctx context.Context, query string) error
	// QueryByDateSku(ctx context.Context, query string) error
	// QueryByService(ctx context.Context, query string) error
	// QueryByDateService(ctx context.Context, query string) error

}

type ProjectDataConfig struct {
	//判断走全部/自定义
	CustomProjectData bool     `json:"custom_project_data"`
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	ProjectIDs        []string `json:"project_ids"`
	ServiceIDs        []string `json:"service_ids"`
	SkusIDs           []string `json:"skus"`
	//
	NegotiatedSavingsEnabled                               bool `json:"negotiated_savings"`
	SavingsProgramsCommittedUsageDiscountEnabled           bool `json:"savings_programs_committed_usage_discount_enable"`
	SavingsProgramsCommittedUsageDiscountDollarBaseEnabled bool `json:"committed_usage_discount_dollar_base_enable"`
	OtherSavingsFreeTierEnabled                            bool `json:"other_savings_free_tier_enable"`
	OtherSavingsPromotionEnabled                           bool `json:"other_savings_promotion_enable"`
	OtherSavingsSustainedUsageDiscountEnabled              bool `json:"other_savings_sustained_usage_discount_enable"`
	OtherSavingsResellerMarginEnabled                      bool `json:"other_savings_reseller_margin_enable"`
	OtherSavingsDiscountEnabled                            bool `json:"other_savings_discount_enable"`
	OtherSavingsSubscriptionBenefitEnabled                 bool `json:"other_savings_subscription_benefit_enable"`
}

func NewProjectDataConfig() *ProjectDataConfig {
	return &ProjectDataConfig{}
}

type ProjectCost struct {
	UsageDate            bigquery.NullDate    `bigquery:"usage_date" json:"usage_date"`                 // DATE 可为 NULL
	ProjectName          bigquery.NullString  `bigquery:"project_name" json:"project_name"`             // STRING 可为 NULL
	ProjectID            bigquery.NullString  `bigquery:"project_id" json:"project_id"`                 // STRING 可为 NULL
	ProjectNumber        bigquery.NullString  `bigquery:"project_number" json:"project_number"`         // STRING 可为 NULL
	InvoiceCost          bigquery.NullFloat64 `bigquery:"invoice_cost" json:"-"`                        // FLOAT 可为 NULL
	InvoiceCostAtListAbs bigquery.NullFloat64 `bigquery:"invoice_cost_at_list_abs" json:"-"`            // FLOAT 可为 NULL
	CostAtList           bigquery.NullFloat64 `bigquery:"cost_at_list" json:"-"`                        // FLOAT 可为 NULL
	UsageCost            bigquery.NullFloat64 `bigquery:"Usage_Cost" json:"Usage_Cost"`                 // FLOAT 可为 NULL
	NegotiatedSavings    bigquery.NullFloat64 `bigquery:"negotiated_savings" json:"negotiated_savings"` // FLOAT 可为 NULL
	SavingsPrograms      bigquery.NullFloat64 `bigquery:"savings_programs" json:"savings_programs"`     // STRING 可为 NULL
	OtherSavings         bigquery.NullFloat64 `bigquery:"other_savings" json:"other_savings"`           // FLOAT 可为 NULL
}
