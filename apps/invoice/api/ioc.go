package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/project"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/services"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/sku"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/gorestful"
	logs "github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&ApiHandler{})
}

type ApiHandler struct {
	ioc.ObjectImpl
	log     *zerolog.Logger
	project invoice.ProjectService
	service invoice.Service
	sku     invoice.SkuService
	// user_binding_roles *mongo.Collection
	// role               *mongo.Collection
	// policy policy.Service
}

func (u *ApiHandler) Init() error {
	u.log = logs.Sub(invoice.AppName)
	u.project = ioc.Controller().Get(project.AppName).(invoice.ProjectService)
	u.service = ioc.Controller().Get(services.AppName).(invoice.Service)
	u.sku = ioc.Controller().Get(sku.AppName).(invoice.SkuService)

	// db := ioc_mongo.DB()
	// u.role = db.Collection("roles")
	// u.stt = ioc.Controller().Get(stt.AppNameV1).(stt.Service)
	// u.policy = ioc.Controller().Get(policy.AppName).(policy.Service)
	// u.user_binding_roles = db.Collection("user_binding_roles")
	u.Registry()
	return nil
}

func (u *ApiHandler) Name() string {
	return invoice.AppName
}

func (u *ApiHandler) Version() string {
	return "v1"
}

// func (i *ApiHandler) Meta() ioc.ObjectMeta {
// 	return ioc.ObjectMeta{
// 		//		CustomPathPrefix: "/s", 必须要/号 http://127.0.0.1:8080/s
// 		CustomPathPrefix: "/",
// 		// CustomPathPrefix: "/s",
// 		Extra: map[string]string{},
// 	}
// }

// ws://localhost:8010/mcenter/api/v1/SpeechToTextV2/ws
func (u *ApiHandler) Registry() {
	tags := []string{"billing invoice month console"}
	ws := gorestful.InitRouter(u)
	ws.Route(ws.POST("/by/data/date-projects").To(u.byDatePojectHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.ProjectDataRequest{}).
		Writes(model.ProjectDateCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", model.ProjectDateCost{}).
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/by/data/projects").To(u.byPojectHandler).
		Doc("基于项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.ProjectDataRequest{}).
		Writes(model.ProjectCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Returns(200, "OK", model.ProjectCost{}).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/all-services-skus").To(u.byAllServicesAllSkusHandler).
		Doc("基于project和指定日期，取所有服务sku").
		Param(ws.QueryParameter("project_id", "项目ID: test-id").DataType("string")).
		Param(ws.QueryParameter("start_date", "开始日期: 2025-xx-xx").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 2025-xx-xx").DataType("string")).
		Reads(model.ProjectDataServiceSkuRequest{}).
		Writes(model.ByDateProjectAllServicesSkusList{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Returns(200, "OK", model.ByDateProjectAllServicesSkusList{}).
		Notes("基于project和指定日期，取所有服务sku"))

	ws.Route(ws.POST("/by/data/date-services").To(u.byDateServiceHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.ServiceDataRequest{}).
		Writes(model.ServiceDateCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", model.ServiceDateCost{}).
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/by/data/services").To(u.byServiceHandler).
		Doc("基于项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.ServiceDataRequest{}).
		Writes(model.ServiceCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Returns(200, "OK", model.ServiceCost{}).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/by/data/date-skus").To(u.byDateSkuHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.SkuDataRequest{}).
		Writes(model.SkuDateCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", model.SkuDateCost{}).
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/by/data/skus").To(u.bySkuHandler).
		Doc("基于项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.SkuDataRequest{}).
		Writes(model.SkuCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Returns(200, "OK", model.SkuCost{}).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.POST("/heru").To(u.byDateSkuHeRuHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("two_decimal_enabled", "是否启用计算精度到两位小数").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 202510").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 202510").DataType("string")).
		Param(ws.QueryParameter("project_ids", "项目ID数组，如: ['myproj-123', 'myproj-456']").DataType("array[string]")).
		Param(ws.QueryParameter("service_ids", "服务ID数组，如: ['6F81-5844-456F']").DataType("array[string]")).
		Param(ws.QueryParameter("skus", "SKU ID数组，如: ['A123-4567-7890']").DataType("array[string]")).
		// 以下为功能开关
		Param(ws.QueryParameter("negotiated_savings", "启用协议价优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_enable", "启用 Resource-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("savings_programs_committed_usage_discount_dollar_base_enable", "启用 Legacy spend-based CUD credits 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_free_tier_enable", "启用 Free Tier 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_promotion_enable", "启用 Promotion 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_sustained_usage_discount_enable", "启用 Sustained Usage Discount").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_reseller_margin_enable", "启用 Reseller Margin 优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_discount_enable", "启用一般折扣优惠计算").DataType("boolean")).
		Param(ws.QueryParameter("other_savings_subscription_benefit_enable", "启用 Subscription Benefit 优惠计算").DataType("boolean")).
		Reads(model.SkuDataRequest{}).
		Writes(model.SkuDateCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", model.SkuDateCost{}).
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Notes("里面包括指定service/sku等"))
}
