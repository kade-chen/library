package api

import (
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/kade-chen/google-billing-console/apps/project"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) streamHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := project.NewProjectDataConfig()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
		return
	}

	// var config project.ProjectDataConfig
	// config.StartDate = "2025-10-01"
	// config.EndDate = "2025-10-02"
	// config.ProjectIDs = []string{"yz-bx3-prod", "kade-poc", "zzshushu", "bw-uat-424708"} // 指定项目
	// config.CustomProjectData = true
	// // config.SkusIDs = []string{"4111-7FF1-D50A"}
	// config.ServiceIDs = []string{"6F81-5844-456A"}
	// config.SkusIDs = []string{"6CB7-B05F-97AD"}
	// 1.默认查询所有的service and skus for project
	if !config.CustomProjectData {
		h.log.Info().Msg("默认查询所有的service and skus for project")
		a, err := h.project.QueryByDateProjectAll(r.Request.Context(), config)
		if err != nil {
			h.log.Info().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, a)
		return
	}
	// 2.自定义查询
	// 2.1 全部service，指定sku
	if len(config.ServiceIDs) == 0 && len(config.SkusIDs) > 0 {
		h.log.Info().Msg("全部service，指定sku")
		a, err := h.project.QueryByDateProjectServicesCustomSku(r.Request.Context(), config)
		if err != nil {
			h.log.Info().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, a)
		return

	}

	// 2.自定义查询
	// 2.2 指定service，全sku
	if len(config.SkusIDs) == 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，全sku")
		a, err := h.project.QueryByDateProjectCustomServicesAllSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Info().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, a)
		return

	}

	// 2.自定义查询
	// 2.3 指定service，指定sku
	if len(config.SkusIDs) > 0 && len(config.ServiceIDs) > 0 {
		h.log.Info().Msg("指定service，指定sku")
		a, err := h.project.QueryByDateProjectCustomServicesCustomSkus(r.Request.Context(), config)
		if err != nil {
			h.log.Info().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
			response.Failed(w, err)
			return
		}
		h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
		response.Success(w, a)
		return

	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
}
