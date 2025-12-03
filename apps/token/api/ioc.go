package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/gorestful"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&tokenHandler{})
}

type tokenHandler struct {
	service token.Service
	log     *zerolog.Logger
	// token.UnimplementedRPCServer
	ioc.ObjectImpl
}

func (h *tokenHandler) Init() error {
	h.log = log.Sub(token.AppName)
	h.service = ioc.Controller().Get(token.AppName).(token.Service)
	h.Registry()
	return nil
}

func (h *tokenHandler) Name() string {
	return token.AppName
}

func (h *tokenHandler) Version() string {
	return "v1"
}

func (h *tokenHandler) Registry() {
	tags := []string{"Token"}
	ws := gorestful.InitRouter(h)

	ws.Route(ws.POST("/").To(h.IssueToken).
		Doc("IssueToken").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Notes("IssueToken").
		Reads(token.IssueTokenRequest{}).
		Returns(200, "ok", token.Token{})) //标签

	ws.Route(ws.DELETE("/").To(h.RevolkToken).
		Doc("RevocationToken").
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Metadata(label.Auth, label.Enable).
		// Metadata(label.PERMISSION_MODE, label.PERMISSION_MODE_ACL.Value()).
		// Metadata(label.Allow, label.AllowAll()).
		Writes(token.Token{}).
		Reads(token.RevolkTokenRequest{}).
		Writes(token.RevolkTokenRequest{}).
		Returns(200, "Ok", token.Token{}).
		Notes("delete the token").
		Returns(404, "Not Found", token.Token{}))

	ws.Route(ws.GET("/verity").To(h.Validate_Token).
		Doc("AuthenticationToken").
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Reads(token.ValicateTokenRequest{}).
		Writes(token.Token{}).
		Notes("verity the token").
		Returns(200, "Ok", token.Token{}))

	// ws.Route(ws.GET("/1").To(func(r *restful.Request, w *restful.Response) {
	// 	w.WriteAsJson(map[string]string{
	// 		"data": "hello mcube1111111",
	// 	})
	// }).Filter(middlewares.NewTokenAuther().Auth_Login))

	// ws.Route(ws.GET("/2").To(func(r *restful.Request, w *restful.Response) {
	// 	// chain.ProcessFilter(req, resp)
	// 	result := r.Request.Context().Value("user").(*user.User)
	// 	fmt.Println("-------", result)
	// 	w.WriteAsJson(map[string]string{
	// 		"data": "hello mcube1111111",
	// 	})
	// }).Filter(middlewares.NewTokenAuther().Auth_Login))

}
