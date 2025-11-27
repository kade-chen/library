package invoice

import "cloud.google.com/go/bigquery"

type InvoiceMonthProjectLabelKeyRequest struct {
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	ProjectIDs []string `json:"project_ids"`
}

func NewInvoiceMonthProjectLabelKeyRequest() *InvoiceMonthProjectLabelKeyRequest {
	return &InvoiceMonthProjectLabelKeyRequest{}
}

type InvoiceMonthProjectLabelKeyList struct {
	LabelKey     bigquery.NullString `bigquery:"keys" json:"keys"`
	LabelValue   bigquery.NullString `bigquery:"value" json:"value"`
	LabelKeyPath bigquery.NullString `bigquery:"key_value_path" json:"key_value_path"`
}

type InvoiceMonthProjectLabelKeyLists struct {
	LabelKey     []string `bigquery:"keys" json:"keys"`
	LabelValue   []string `bigquery:"value" json:"value"`
	LabelKeyPath []string `bigquery:"key_value_path" json:"key_value_path"`
}
