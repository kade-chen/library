package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
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
	tags := []string{"Speech To Text V1 Client"}
	ws := gorestful.InitRouter(u)
	ws.Route(ws.POST("/").To(u.streamHandler).
		Doc("Websocket Streaming Recognize").
		Reads(project.ProjectDataConfig{}).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Notes("Websocket Streaming Recognize"))
	// ws.Route(ws.GET("/ws").To(u.streamingRecognize).
	// 	Doc("Websocket Streaming Recognize").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags). //标签
	// 	Notes("Websocket Streaming Recognize").
	// 	Filter(middlewares.NewTokenAuther().Auth_Login))

	// Writes(user.User{}).
	// Reads(user.CreateUserRequest{}).
	// Returns(200, "Ok", user.User{}).
	// )

	// LocalStreamingRecognize
	// ws.Route(ws.GET("/").To(u.localStreamingRecognize).
	// 	Doc("create user").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags). //标签
	// 	Notes("create user"))
	// Writes(user.User{}).
	// Reads(user.CreateUserRequest{}).
	// Returns(200, "Ok", user.User{}).
	// )
}
