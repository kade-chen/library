package organization

import (
	"strings"

	"cloud.google.com/go/bigquery"
)

// new QueryOrganizationRequest request
func NewListOrganizationRequest(req *ListOrganizationRequest) *ListOrganizationRequest {
	if req == nil {
		return &ListOrganizationRequest{
			// Page:         NewPageRequest(20, 1),
			Page: &PageRequest{},
		}
	}
	if req.Page == nil {
		req.Page = &PageRequest{}
	}
	return req
}

func NewOrganizationSet() *OrganizationSet {
	return &OrganizationSet{
		Items: []*Organization{},
	}
}

func (u *OrganizationSet) Add(item *Organization) {
	u.Items = append(u.Items, item)
}

func (r *ListOrganizationRequest) PageSQL() (string, []bigquery.QueryParameter) {
	conditions := []string{}
	params := []bigquery.QueryParameter{}

	// LIMIT
	if r.Page.PageSize > 0 {
		conditions = append(conditions, " LIMIT @limit")
		params = append(params,
			bigquery.QueryParameter{Name: "limit", Value: int64(r.Page.PageSize)},
		)
	}

	// OFFSET
	if r.Page.Offset > 0 {
		conditions = append(conditions, " OFFSET @offset")
		params = append(params,
			bigquery.QueryParameter{Name: "offset", Value: int64(r.Page.Offset)},
		)
	}
	return strings.Join(conditions, " "), params
}

func (r *ListOrganizationRequest) WhereSQL() (string, []bigquery.QueryParameter) {
	conditions := []string{}
	params := []bigquery.QueryParameter{}

	// Organization
	if len(r.Names) > 0 {
		conditions = append(conditions, `spec.organization_detail.sub_organization IN UNNEST(@Organization)`)
		params = append(params, bigquery.QueryParameter{Name: "Organization", Value: r.Names})
	}

	if len(conditions) == 0 {
		return "", params
	}

	return "WHERE " + strings.Join(conditions, " AND "), params
}
