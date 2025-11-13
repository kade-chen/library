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

	// 1.默认查询所有的service and skus for project
	if !config.CustomProjectData {
		h.log.Info().Msg("默认查询所有的service and skus for project")
		projectCost, err := h.project.QueryByDateProjectAll(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return
	}
	// 2.自定义查询
	// 2.1 全部service，指定sku
	if len(config.ServiceIDs) == 0 && len(config.SkusIDs) > 0 {
		h.log.Info().Msg("全部service，指定sku")
		projectCost, err := h.project.QueryByDateProjectServicesCustomSku(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return

	}

	// 2.自定义查询
	// 2.2 指定service，全sku
	if len(config.SkusIDs) == 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，全sku")
		projectCost, err := h.project.QueryByDateProjectCustomServicesAllSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return

	}

	// 2.自定义查询
	// 2.3 指定service，指定sku
	if len(config.SkusIDs) > 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，指定sku")
		projectCost, err := h.project.QueryByDateProjectCustomServicesCustomSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewProjectDataConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	// 1.默认查询所有的service and skus for project
	if !config.CustomProjectData {
		h.log.Info().Msg("默认查询所有的service and skus for project")
		projectCost, err := h.project.QueryByProjectAll(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return
	}
	// 2.自定义查询
	// 2.1 全部service，指定sku
	if len(config.ServiceIDs) == 0 && len(config.SkusIDs) > 0 {
		h.log.Info().Msg("全部service，指定sku")
		projectCost, err := h.project.QueryByProjectServicesCustomSku(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return

	}

	// 2.自定义查询
	// 2.2 指定service，全sku
	if len(config.SkusIDs) == 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，全sku")
		projectCost, err := h.project.QueryByProjectCustomServicesAllSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return

	}

	// 2.自定义查询
	// 2.3 指定service，指定sku
	if len(config.SkusIDs) > 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，指定sku")
		projectCost, err := h.project.QueryByProjectCustomServicesCustomSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, projectCost)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
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
