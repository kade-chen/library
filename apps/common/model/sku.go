package model

import "cloud.google.com/go/bigquery"

type SkusList struct {
	ServiceID      bigquery.NullString `bigquery:"service_id" json:"service_id"`
	SKUId          bigquery.NullString `bigquery:"sku_id" json:"sku_id"`
	SKUDesc        bigquery.NullString `bigquery:"sku_desc" json:"sku_desc"`
	ServicePath    bigquery.NullString `bigquery:"service_path" json:"service_path"`
	ServiceSKUPath bigquery.NullString `bigquery:"service_sku_path" json:"service_sku_path"`
}
