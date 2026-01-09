package api

import (
	"github.com/emicklei/go-restful/v3"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
	"github.com/kade-chen/google-billing-console/tools/trances"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) byDatePojectHandler(r *restful.Request, w *restful.Response) {
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDatePojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}
	//3.调用每天项目费用接口
	projectCost, err := h.project.QueryByDateProject(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDatePojectAPI ✅", trancesID)
	response.Success(w, projectCost)
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByPojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	projectCost, err := h.project.QueryByProject(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByPojectAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateServiceHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDateServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	projectCost, err := h.service.QueryByDateService(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDateServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byServiceHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	projectCost, err := h.service.QueryByService(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateSkuHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDateSkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	projectCost, err := h.sku.QueryByDateSku(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDateSkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) bySkuHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateBySkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	projectCost, err := h.sku.QueryBySku(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateBySkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byAllServicesAllSkusHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByServicesSkusAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataServiceSkusRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	a, err := h.project.QueryByDateProjectAllServicesAllSkus(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByServicesSkusAPI ✅", trancesID)
	response.Success(w, a)
}

func (h *ApiHandler) byAllLabelKeyHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	trancesID := trances.NewTraceID()

	// 注入 trances_id 到 context
	r.Request = trances.NewTraceIDToRequest(r.Request, trancesID)

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByLabelKeyAPI", trancesID)

	//2.read the request body parametars
	config := model.NewUsageDateProjectLabelKeyRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err))
		return
	}

	a, err := h.labelkey.QueryByUsageDatProjectLabelKeyAll(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByLabelKeyAPI ✅", trancesID)
	response.Success(w, a)
}
