package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/application"
	"github.com/go-openapi/spec"
)

// ---------------------------------------------------------------
// Start 启动服务
func (h *Http) Start(ctx context.Context) {
	h.log.Info().Msgf("HTTP服务启动成功, 监听地址: http://%s", h.Addr())
	if err := h.server.ListenAndServe(); err != nil {
		h.log.Error().Msg(err.Error())
	}
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

// ---------------------------------------------------------------

// Stop 停止server
func (h *Http) Stop(ctx context.Context) error {
	h.log.Info().Msg("start graceful shutdown")
	// 优雅关闭HTTP服务
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http graceful shutdown timeout, force exit")
	}
	return nil
}

func (h *Http) ApiObjectPathPrefix(obj ioc.Object) string {
	cp := obj.Meta().CustomPathPrefix
	if cp != "" {
		// h.log.Debug().Msgf("use custom path prefix : http://127.0.0.1:8080%s", cp)
		return cp
	}
	// h.log.Debug().Msgf("use custom path prefix : http://127.0.0.1:8080%s/%s/%s", h.HTTPPrefix(), obj.Version(), obj.Name())
	return fmt.Sprintf("%s/%s/%s", h.HTTPPrefix(), obj.Version(), obj.Name())
}

func (h *Http) HTTPPrefix() string {
	u, err := url.JoinPath("/"+application.Get().AppName, h.PathPrefix)
	if err != nil {
		return fmt.Sprintf("/%s/%s", application.Get().AppName, h.PathPrefix)
	}
	return u
}

func (h *Http) ApiObjectAddr(obj ioc.Object) string {
	return fmt.Sprintf("http://%s%s", h.Addr(), h.ApiObjectPathPrefix(obj))
}

func (h *Http) SetRouter(r http.Handler) {
	h.router = r
}

func (h *Http) IsEnable() bool {
	if h.Enable == nil {
		return h.router != nil
	}

	return *h.Enable
}

func (a *Http) SwagerDocs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       application.Get().AppName,
			Description: application.Get().AppDescription,
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: application.Short(),
		},
	}
}
