package query

import (
	"net/url"
	"strconv"
)

//nolint:tagliatelle
type ListCompaniesRequest struct {
	Query           string   `json:"query"`
	BusinessTypeIDs []string `json:"business_type_ids"`
	Limit           int      `json:"limit"`
}

func NewListCompaniesRequestFromQueryParams(values url.Values) *ListCompaniesRequest {
	limit, _ := strconv.Atoi(values["limit"][0])

	return &ListCompaniesRequest{
		Query:           values.Get("q"),
		BusinessTypeIDs: values["business_type_ids"],
		Limit:           limit,
	}
}
