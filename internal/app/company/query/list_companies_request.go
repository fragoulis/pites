package query

import "net/url"

//nolint:tagliatelle
type ListCompaniesRequest struct {
	Query           string   `json:"query"`
	BusinessTypeIDs []string `json:"business_type_ids"`
}

func NewListCompaniesRequestFromQueryParams(values url.Values) *ListCompaniesRequest {
	return &ListCompaniesRequest{
		Query:           values.Get("q"),
		BusinessTypeIDs: values["business_type_ids"],
	}
}
