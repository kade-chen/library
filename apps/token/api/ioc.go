package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/kade-chen/google-billing-console/apps/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/token"
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
	jwt     auth.Service
	// token.UnimplementedRPCServer
	ioc.ObjectImpl
}

func (h *tokenHandler) Init() error {
	h.log = log.Sub(token.AppName)
	h.service = ioc.Controller().Get(token.AppName).(token.Service)
	h.jwt = ioc.Controller().Get(auth.AppName).(auth.Service)
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
	ws.Route(ws.POST("/login").To(h.IssueToken).
		Doc("IssueToken").
		Param(ws.QueryParameter("username", "Username, example: example@example.com").DataType("string")).
		Param(ws.QueryParameter("password", "Password, example: xxxxxxxx").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Notes("IssueToken").
		Reads(token.IssueTokenRequest{}).
		Returns(200, "ok", token.Token{}).
		Returns(404, "Not Found", token.Token{})) //标签

	ws.Route(ws.DELETE("/logout").To(h.RevolkToken).
		Doc("RevolkToken").
		Param(ws.QueryParameter("Authorization", "Authorization, example: Bearer xxxxxxxx").DataType("string")).
		// Param(ws.QueryParameter("password", "Password, example: xxxxxxxx").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		// Metadata(label.Auth, label.Enable).
		// Metadata(label.Allow, label.AllowAll()).
		Writes(token.Token{}).
		Reads(token.RevolkTokenRequest{}).
		Writes(token.RevolkTokenRequest{}).
		Returns(404, "Not Found", token.Token{}).
		Notes("delete the token").
		Filter(h.jwt.Auth).
		Returns(200, "Ok", token.Token{}))

	ws.Route(ws.POST("/verity").To(h.Validate_Token).
		Doc("ValidateToken").
		Param(ws.QueryParameter("access_token", "jwrToken, example: Bearer xxxxx").DataType("string")).
		Param(ws.QueryParameter("access_token_name", "暂未开放使用").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Reads(token.ValicateTokenRequest{}).
		Writes(model.ValidateToken{}).
		Notes("verity the token").
		// Filter(h.jwt.Auth).
		Notes("验证token是否过期").
		Returns(200, "Ok", model.ValidateToken{}))

	// ws.Route(ws.GET("/test").To(h.test).
	// 	Doc("AuthenticationToken").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags). //标签
	// 	Reads(token.ValicateTokenRequest{}).
	// 	Writes(token.Token{}).
	// 	Notes("verity the token").
	// 	Filter(h.jwt.Auth).
	// 	// Filter(middlewares.NewTokenAuthMiddleware().Auth).
	// 	Returns(200, "Ok", token.Token{}))

	ws.Route(ws.POST("/refresh").To(h.Refresh_Token).
		Doc("RefreshToken").
		Metadata(restfulspec.KeyOpenAPITags, tags). //标签
		Reads(token.ValicateTokenRequest{}).
		Writes(token.Token{}).
		Notes("Refresh the token").
		Consumes(
			restful.MIME_JSON,
			"application/json; charset=utf-8",
			"",    // 没有 Content-Type
			"*/*", // ⭐ 兜底（关键）
		).
		Produces(restful.MIME_JSON).
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
