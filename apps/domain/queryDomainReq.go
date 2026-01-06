package domain

import (
	"strings"

	"cloud.google.com/go/bigquery"
)

// new QueryDomainRequest request
func NewListDomainRequest(req *ListDomainRequest) *ListDomainRequest {
	if req == nil {
		return &ListDomainRequest{
			// Page:         NewPageRequest(20, 1),
			Page: &PageRequest{},
		}
	}
	if req.Page == nil {
		req.Page = &PageRequest{}
	}
	return req
}

func NewDomainSet() *DomainSet {
	return &DomainSet{
		Items: []*Domain{},
	}
}

func (u *DomainSet) Add(item *Domain) {
	u.Items = append(u.Items, item)
}

func (r *ListDomainRequest) PageSQL() (string, []bigquery.QueryParameter) {
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

func (r *ListDomainRequest) WhereSQL() (string, []bigquery.QueryParameter) {
	conditions := []string{}
	params := []bigquery.QueryParameter{}

	// domain
	if len(r.Names) > 0 {
		conditions = append(conditions, `spec.name IN UNNEST(@domain)`)
		params = append(params, bigquery.QueryParameter{Name: "domain", Value: r.Names})
	}

	if len(conditions) == 0 {
		return "", params
	}

	return "WHERE " + strings.Join(conditions, " AND "), params
}
