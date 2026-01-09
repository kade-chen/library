package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/kade-chen/google-billing-console/apps/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
	"github.com/kade-chen/google-billing-console/apps/organization"
	"github.com/kade-chen/google-billing-console/apps/usagedate"
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
	log    *zerolog.Logger
	jwt    auth.Service
	domian organization.Service
	// user_binding_roles *mongo.Collection
	// role               *mongo.Collection
	// policy policy.Service
}

func (u *ApiHandler) Init() error {
	u.log = logs.Sub(organization.AppName)
	u.log.Debug().Msgf("---------%s API begin init.......---------", usagedate.AppName)
	u.jwt = ioc.Controller().Get(auth.AppName).(auth.Service)
	// u.role = db.Collection("roles")
	u.domian = ioc.Controller().Get(organization.AppName).(organization.Service)
	// u.policy = ioc.Controller().Get(policy.AppName).(policy.Service)
	// u.user_binding_roles = db.Collection("user_binding_roles")
	u.Registry()
	u.log.Debug().Msgf("---------%s All API init succeessful✅---------", usagedate.AppName)
	return nil
}

func (u *ApiHandler) Name() string {
	return organization.AppName
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
	tags := []string{"Organizations"}
	ws := gorestful.InitRouter(u)
	ws.Route(ws.GET("/list").To(u.listOrginzationsHandler).
		Doc("基于日期的项目费用统计").
		Param(ws.QueryParameter("Authorization", "Authorization, example: Bearer xxxxxxxx").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Filter(u.jwt.Auth).
		Returns(200, "OK", model.ProjectDateCost{}).
		// Filter(middlewares.NewTokenAuther().Auth_Login).
		Notes("列出jwt_token所对应用户的Organizations"))
}
