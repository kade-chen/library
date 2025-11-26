package api

import (
	"fmt"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	model "github.com/kade-chen/google-billing-console/apps/common/model/invoice"
	"github.com/kade-chen/google-billing-console/tools/csv"
	"github.com/kade-chen/google-billing-console/tools/validate"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) byDatePojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewProjectDataRequest()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
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

	projectCost, err := h.project.QueryByDateProject(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	err = csv.WriteStructToCSV("output.csv", projectCost)
	if err != nil {
		h.log.Error().Msgf("shengchengcsv fiald %v", err)
		response.Failed(w, err)
		return
	}

	fmt.Println("写入完成")
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
	response.Success(w, projectCost)
	// return
}

func (h *ApiHandler) byPojectHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewProjectRequest()
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

func (h *ApiHandler) byDateServiceHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewServiceDataRequest()
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
	config := model.NewServiceRequest()
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
	config := model.NewSkuDataRequest()
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
	config := model.NewSkuRequest()
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

func (h *ApiHandler) byAllServicesAllSkusHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))
	//2.read the request body parametars
	config := model.NewProjectDataServiceSkuRequest()
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

func (h *ApiHandler) byDateSkuHeRuHandler(r *restful.Request, w *restful.Response) {
	// 生成唯一请求ID
	reqID := uuid.New().String()
	h.log.Info().Msgf("request_id=%s time=%s 接口被调用", reqID, time.Now().Format(time.RFC3339))

	//2.read the request body parametars
	config := model.NewSkuDataRequest()
	if err := r.ReadEntity(config); err != nil {
		response.Failed(w, exception.NewInternalServerError("read request struct: %v; error: %v", config, err))
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
	projectCost, err := h.sku.QueryByDateSkuHeru(r.Request.Context(), config)
	if err != nil {
		h.log.Error().Msgf("request_id=%s time=%s 接口调用失败✅", reqID, time.Now().Format(time.RFC3339))
		response.Failed(w, err)
		return
	}
	h.log.Info().Msgf("request_id=%s time=%s 接口已完成✅", reqID, time.Now().Format(time.RFC3339))
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
