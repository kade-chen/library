package api

import (
	"context"
	"fmt"

	"github.com/emicklei/go-restful/v3"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	"github.com/kade-chen/google-billing-console/tools/csv"
	"github.com/kade-chen/google-billing-console/tools/validate"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) byDatePojectHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByDatePojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.project.QueryByDateProject(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	// err = csv.WriteStructToCSV("output.csv", projectCost)
	// if err != nil {
	// 	h.log.Error().Msgf("shengchengcsv fiald %v", err)
	// 	response.Failed(w, err)
	// 	return
	// }

	// fmt.Println("写入完成")
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByDatePojectAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByPojectAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.project.QueryByProject(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByPojectAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateServiceHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByDateServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.service.QueryByDateService(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByDateServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byServiceHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByServiceAPI", trancesID)

	//2.read the request body parametars
	config := model.NewServiceRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.service.QueryByService(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	// err = csv.WriteStructToCSV("output.csv", projectCost)
	// if err != nil {
	// 	h.log.Error().Msgf("shengchengcsv fiald %v", err)
	// 	response.Failed(w, err)
	// 	return
	// }

	fmt.Println("写入完成")
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByServiceAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byDateSkuHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByDateSkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.sku.QueryByDateSku(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByDateSkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) bySkuHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthBySkuAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}

	projectCost, err := h.sku.QueryBySku(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthBySkuAPI ✅", trancesID)
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byAllServicesAllSkusHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByLabelKeyAPI", trancesID)

	//2.read the request body parametars
	config := model.NewProjectDataServiceSkuRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	a, err := h.project.QueryByDateProjectAllServicesAllSkus(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByLabelKeyAPI ✅", trancesID)
	response.Success(w, a)
}

func (h *ApiHandler) byInvoiceMonthLabelKeyHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByLabelKeyAPI", trancesID)

	//2.read the request body parametars
	config := model.NewInvoiceMonthProjectLabelKeyRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	a, err := h.labelkey.QueryByInvoiceMonthProjectLabelKeyAll(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByLabelKeyAPI ✅", trancesID)
	response.Success(w, a)
}

func (h *ApiHandler) byDateSkuHeRuHandler(r *restful.Request, w *restful.Response) {
	ctx := context.WithValue(r.Request.Context(), "claims", r.Attribute("claims").(*authModel.TokenAuthMiddleware))
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID

	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface InvoiceMonthByServicesSkusAPI", trancesID)

	//2.read the request body parametars
	config := model.NewSkuDataRequest()
	if err := r.ReadEntity(&config); err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, exception.NewInternalServerError("trances_id=%s, ERROR: %v", trancesID, err))
		return
	}

	if err := validate.ValidateYYYYMM(config.StartDate); err != nil {
		response.Failed(w, err)
		return
	}

	if err := validate.ValidateYYYYMM(config.EndDate); err != nil {
		response.Failed(w, err)
		return
	}
	projectCost, err := h.sku.QueryByDateSkuHeru(ctx, config)
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", trancesID, err)
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for InvoiceMonthByServicesSkusAPI ✅", trancesID)

	err = csv.WriteStructToCSV("output.csv", projectCost)
	if err != nil {
		h.log.Error().Msgf("shengchengcsv fiald %v", err)
		response.Failed(w, err)
		return
	}

	fmt.Println("写入完成")
	response.Success(w, projectCost)
	// return
}
