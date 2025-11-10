package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/google-billing-console/apps/project"
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
	project project.Service
	// user_binding_roles *mongo.Collection
	// role               *mongo.Collection
	// policy policy.Service
}

func (u *ApiHandler) Init() error {
	u.log = logs.Sub(project.AppName)
	u.project = ioc.Controller().Get(project.AppName).(project.Service)

	// db := ioc_mongo.DB()
	// u.role = db.Collection("roles")
	// u.stt = ioc.Controller().Get(stt.AppNameV1).(stt.Service)
	// u.policy = ioc.Controller().Get(policy.AppName).(policy.Service)
	// u.user_binding_roles = db.Collection("user_binding_roles")
	u.Registry()
	return nil
}

func (u *ApiHandler) Name() string {
	return project.AppName
}

func (u *ApiHandler) Version() string {
	return "v1"
}

func (i *ApiHandler) Meta() ioc.ObjectMeta {
	return ioc.ObjectMeta{
		//		CustomPathPrefix: "/s", 必须要/号 http://127.0.0.1:8080/s
		CustomPathPrefix: "/billing-console",
		// CustomPathPrefix: "/s",
		Extra: map[string]string{},
	}
}

// ws://localhost:8010/mcenter/api/v1/SpeechToTextV2/ws
func (u *ApiHandler) Registry() {
	tags := []string{"billing console"}
	ws := gorestful.InitRouter(u)
	ws.Route(ws.POST("/by/data/project").To(u.streamHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("custom_project_data", "是否启用自定义项目: true/false 默认是查询全部").DataType("boolean")).
		Param(ws.QueryParameter("start_date", "开始日期: 2025-11-01").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 2025-11-05").DataType("string")).
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
		Reads(project.ProjectDataConfig{}).
		Writes(project.ProjectCost{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", project.ProjectCost{}).
		Notes("里面包括指定service/sku等"))

	ws.Route(ws.GET("/all-services-skus").To(u.streamHandler1).
		Doc("基于project和指定日期，取所有服务sku").
		Param(ws.QueryParameter("project_id", "项目ID: test-id").DataType("string")).
		Param(ws.QueryParameter("start_date", "开始日期: 2025-xx-xx").DataType("string")).
		Param(ws.QueryParameter("end_date", "结束日期: 2025-xx-xx").DataType("string")).
		Reads(model.ProjectDataRequest{}).
		Writes(model.ByDateProjectAllServicesSkusList{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Returns(200, "OK", model.ByDateProjectAllServicesSkusList{}).
		Notes("基于project和指定日期，取所有服务sku"))
}
