package api

import (
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/kade-chen/google-billing-console/apps/common/model"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) byDatePojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewProjectDataConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.project.QueryByDateProject(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewProjectConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.project.QueryByProject(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byAllServicesAllSkusHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))
	//2.read the request body parametars
	config := model.NewProjectDataRequest()
	if err := r.ReadEntity(config); err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; ERROR: %v", config, err))
		return
	}

	a, err := h.project.QueryByDateProjectAllServicesAllSkus(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, a)
}



func (h *ApiHandler) byDateServiceHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewServiceDataConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.service.QueryByDateService(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byServiceHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewServiceConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.service.QueryByService(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}



func (h *ApiHandler) byDateSkuHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewSkuDataConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.sku.QueryByDateSku(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) bySkuHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewSkuConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	projectCost, err := h.sku.QueryBySku(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}
