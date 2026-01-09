package api

import (
	"github.com/emicklei/go-restful/v3"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) listOrginzationsHandler(r *restful.Request, w *restful.Response) {

	domains := r.Attribute("claims").(*authModel.TokenAuthMiddleware).Domains
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID
	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface ListOrginzationtsAPI", trancesID)

	//3.调用每天项目费用接口
	domainSet, err := h.domian.ListDoamin(r.Request.Context(), &domain.ListDomainRequest{
		Names: domains,
	})
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	if domainSet.Total != int64(len(domains)) {
		var cc []string
		for _, v := range domainSet.Items {
			cc = append(cc, v.Spec.Name)
		}
		h.log.Warn().Msgf("trances_id=%s, ERROR: jwt_token Organizational privilege refresh or forgery", r.Request.Context().Value("trances_id"))
		response.Success(w, cc)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for ListOrginzationtsAPI ✅", trancesID)
	response.Success(w, domains)
}
