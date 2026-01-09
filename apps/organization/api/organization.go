package api

import (
	"github.com/emicklei/go-restful/v3"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/organization"
	"github.com/kade-chen/library/http/response"
)

func (h *ApiHandler) listOrginzationsHandler(r *restful.Request, w *restful.Response) {

	organizations := r.Attribute("claims").(*authModel.TokenAuthMiddleware).Organizations
	trancesID := r.Attribute("claims").(*authModel.TokenAuthMiddleware).TrancesID
	h.log.Info().Msgf("trances_id=%s, The User begins calling the interface ListOrginzationtsAPI", trancesID)

	//3.调用每天项目费用接口
	organizationSet, err := h.domian.ListOrganizations(r.Request.Context(), &organization.ListOrganizationRequest{
		Names: organizations,
	})
	if err != nil {
		h.log.Error().Msgf("trances_id=%s, ERROR: %v", r.Request.Context().Value("trances_id"), err)
		response.Failed(w, err)
		return
	}
	if organizationSet.Total != int64(len(organizations)) {
		var cc []string
		for _, v := range organizationSet.Items {
			cc = append(cc, v.Spec.SubOrganization)
		}
		h.log.Warn().Msgf("trances_id=%s, ERROR: jwt_token Organizational privilege refresh or forgery", r.Request.Context().Value("trances_id"))
		response.Success(w, cc)
		return
	}
	h.log.Info().Msgf("trances_id=%s, The User calling the interface Successful for ListOrginzationtsAPI ✅", trancesID)
	response.Success(w, organizations)
}
