package model

// ✅ 聚合结构体
type ByDateProjectAllServicesSkusList struct {
	Services []ServicesList `json:"services"`
	Skus     []SkusList     `json:"skus"`
}
