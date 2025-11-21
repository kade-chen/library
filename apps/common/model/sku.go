package model

import "cloud.google.com/go/bigquery"

type SkusList struct {
	ServiceID      bigquery.NullString `bigquery:"service_id" json:"service_id"`
	SKUId          bigquery.NullString `bigquery:"sku_id" json:"sku_id"`
	SKUDescribe    bigquery.NullString `bigquery:"sku_describe" json:"sku_describe"`
	ServicePath    bigquery.NullString `bigquery:"service_path" json:"service_path"`
	ServiceSKUPath bigquery.NullString `bigquery:"service_sku_path" json:"service_sku_path"`
}

type SkuDataConfig struct {
	TwoDecimalEnabled bool `json:"two_decimal_enabled"`
	//判断走全部/自定义
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	ProjectIDs []string `json:"project_ids"`
	ServiceIDs []string `json:"service_ids"`
	SkusIDs    []string `json:"skus"`
	//
	NegotiatedSavingsEnabled                               bool `json:"negotiated_savings"`
	SavingsProgramsCommittedUsageDiscountEnabled           bool `json:"savings_programs_committed_usage_discount_enable"`
	SavingsProgramsCommittedUsageDiscountDollarBaseEnabled bool `json:"savings_programs_committed_usage_discount_dollar_base_enable"`
	OtherSavingsFreeTierEnabled                            bool `json:"other_savings_free_tier_enable"`
	OtherSavingsPromotionEnabled                           bool `json:"other_savings_promotion_enable"`
	OtherSavingsSustainedUsageDiscountEnabled              bool `json:"other_savings_sustained_usage_discount_enable"`
	OtherSavingsResellerMarginEnabled                      bool `json:"other_savings_reseller_margin_enable"`
	OtherSavingsDiscountEnabled                            bool `json:"other_savings_discount_enable"`
	OtherSavingsSubscriptionBenefitEnabled                 bool `json:"other_savings_subscription_benefit_enable"`
}

func NewSkuDataConfig() *SkuDataConfig {
	return &SkuDataConfig{}
}

type SkuDateCost struct {
	UsageDate          bigquery.NullDate   `bigquery:"usage_date" json:"usage_date"`                   // DATE 可为 NULL
	SkuDescription     bigquery.NullString `bigquery:"sku_description" json:"sku_description"`         // STRING 可为 NULL
	ServiceDescription bigquery.NullString `bigquery:"service_description" json:"service_description"` // STRING 可为 NULL
	SkuId              bigquery.NullString `bigquery:"sku_id" json:"sku_id"`                           // STRING 可为 NULL
	Usage              bigquery.NullString `bigquery:"usage" json:"usage"`                             // FLOAT 可为 NULL
	UsageDetails       bigquery.NullString `bigquery:"usage_details" json:"usage_details"`             // STRING 可为 NULL
	// InvoiceCost          bigquery.NullFloat64 `bigquery:"invoice_cost" json:"-"`                          // FLOAT 可为 NULL
	// InvoiceCostAtListAbs bigquery.NullFloat64 `bigquery:"invoice_cost_at_list_abs" json:"-"`              // FLOAT 可为 NULL
	// CostAtList           bigquery.NullFloat64 `bigquery:"cost_at_list" json:"-"`                          // FLOAT 可为 NULL
	UsageCost         bigquery.NullFloat64 `bigquery:"Usage_Cost" json:"Usage_Cost"`                 // FLOAT 可为 NULL
	NegotiatedSavings bigquery.NullFloat64 `bigquery:"negotiated_savings" json:"negotiated_savings"` // FLOAT 可为 NULL
	SavingsPrograms   bigquery.NullFloat64 `bigquery:"savings_programs" json:"savings_programs"`     // STRING 可为 NULL
	OtherSavings      bigquery.NullFloat64 `bigquery:"other_savings" json:"other_savings"`           // FLOAT 可为 NULL
	SubTotal          bigquery.NullFloat64 `bigquery:"sub_total" json:"sub_total"`                   // FLOAT 可为 NULL
}

type SkuConfig struct {
	TwoDecimalEnabled bool     `json:"two_decimal_enabled"`
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	ProjectIDs        []string `json:"project_ids"`
	ServiceIDs        []string `json:"service_ids"`
	SkusIDs           []string `json:"skus"`
	//
	NegotiatedSavingsEnabled                               bool `json:"negotiated_savings"`
	SavingsProgramsCommittedUsageDiscountEnabled           bool `json:"savings_programs_committed_usage_discount_enable"`
	SavingsProgramsCommittedUsageDiscountDollarBaseEnabled bool `json:"savings_programs_committed_usage_discount_dollar_base_enable"`
	OtherSavingsFreeTierEnabled                            bool `json:"other_savings_free_tier_enable"`
	OtherSavingsPromotionEnabled                           bool `json:"other_savings_promotion_enable"`
	OtherSavingsSustainedUsageDiscountEnabled              bool `json:"other_savings_sustained_usage_discount_enable"`
	OtherSavingsResellerMarginEnabled                      bool `json:"other_savings_reseller_margin_enable"`
	OtherSavingsDiscountEnabled                            bool `json:"other_savings_discount_enable"`
	OtherSavingsSubscriptionBenefitEnabled                 bool `json:"other_savings_subscription_benefit_enable"`
}

func NewSkuConfig() *SkuConfig {
	return &SkuConfig{}
}

type SkuCost struct {
	SkuDescription     bigquery.NullString `bigquery:"sku_description" json:"sku_description"`         // STRING 可为 NULL
	ServiceDescription bigquery.NullString `bigquery:"service_description" json:"service_description"` // STRING 可为 NULL
	SkuId              bigquery.NullString `bigquery:"sku_id" json:"sku_id"`                           // STRING 可为 NULL
	Usage              bigquery.NullString `bigquery:"usage" json:"usage"`                             // STRING 可为 NULL
	UsageDetails       bigquery.NullString `bigquery:"usage_details" json:"usage_details"`             // STRING 可为 NULL
	// InvoiceCost          bigquery.NullFloat64 `bigquery:"invoice_cost" json:"-"`                          // FLOAT 可为 NULL
	// InvoiceCostAtListAbs bigquery.NullFloat64 `bigquery:"invoice_cost_at_list_abs" json:"-"`              // FLOAT 可为 NULL
	// CostAtList           bigquery.NullFloat64 `bigquery:"cost_at_list" json:"-"`                          // FLOAT 可为 NULL
	UsageCost         bigquery.NullFloat64 `bigquery:"Usage_Cost" json:"Usage_Cost"`                 // FLOAT 可为 NULL
	NegotiatedSavings bigquery.NullFloat64 `bigquery:"negotiated_savings" json:"negotiated_savings"` // FLOAT 可为 NULL
	SavingsPrograms   bigquery.NullFloat64 `bigquery:"savings_programs" json:"savings_programs"`     // STRING 可为 NULL
	OtherSavings      bigquery.NullFloat64 `bigquery:"other_savings" json:"other_savings"`           // FLOAT 可为 NULL
	SubTotal          bigquery.NullFloat64 `bigquery:"sub_total" json:"sub_total"`                   // FLOAT 可为 NULL
	ChangeRate        bigquery.NullString  `bigquery:"change_rate" json:"change_rate"`               // FLOAT 可为 NULL
}
