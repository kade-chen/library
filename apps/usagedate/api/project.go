package api

import (
	"context"
	"fmt"

	"github.com/emicklei/go-restful/v3"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/usagedate"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) byDatePojectHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDatePojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}
	//3.调用每天项目费用接口
	projectCost, err := h.project.QueryByDateProject(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDatePojectAPI ✅", trancesID)
	response.Success(w, projectCost)
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByPojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	projectCost, err := h.project.QueryByProject(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByPojectAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateServiceHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDateServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	projectCost, err := h.service.QueryByDateService(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDateServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byServiceHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	projectCost, err := h.service.QueryByService(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateSkuHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByDateSkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	projectCost, err := h.sku.QueryByDateSku(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByDateSkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) bySkuHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateBySkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	projectCost, err := h.sku.QueryBySku(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateBySkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byAllServicesAllSkusHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID
	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByServicesSkusAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataServiceSkusRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	fmt.Println("-----", config.OrganizationBqTable)
	a, err := h.project.QueryByDateProjectAllServicesAllSkus(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByServicesSkusAPI ✅", trancesID)
	response.Success(w, a)
}

func (h *ApiHandler) byAllLabelKeyHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface UsageDateByLabelKeyAPI", trancesID)

	//2.read the request body parametars
	config := model.NewUsageDateProjectLabelKeyRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	a, err := h.labelkey.QueryByUsageDatProjectLabelKeyAll(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for UsageDateByLabelKeyAPI ✅", trancesID)
	response.Success(w, a)
}
