package sku

import "context"

const (
	AppName = "skus"
)

type Service interface {
	
	QueryByDateProjectSKUsAll(ctx context.Context, query string) ([]SkusList, error)
}

type SkusList struct {
	ServiceID      string `json:"service_id"`
	SKUId          string `json:"sku_id"`
	SKUDesc        string `json:"sku_describe"`
	ServicePath    string `json:"service_path"`
	ServiceSKUPath string `json:"service_sku_path"`
}
