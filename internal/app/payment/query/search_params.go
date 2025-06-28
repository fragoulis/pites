package query

import (
	"net/url"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/utils"
)

const (
	DefaultListCount = int64(50)
)

//nolint:tagliatelle
type ListPaymentsRequest struct {
	Query        string   `json:"q"`
	PaymentIDs   []string `json:"member_ids"`
	UserIDs      []string `json:"user_ids"`
	WithReceipts *bool    `json:"with_receipts"`
	Page         int64    `json:"page"`
	PerPage      int64    `json:"per_page"`
	Sort         string   `json:"sort"`
}

func NewListPaymentsRequest(values url.Values) *ListPaymentsRequest {
	page := 0
	perPage := 0

	var withReceipts *bool

	pageRaw := values.Get("page")
	if pageRaw != "" {
		page, _ = strconv.Atoi(pageRaw)
	}

	perPageRaw := values.Get("per_page")
	if perPageRaw != "" {
		perPage, _ = strconv.Atoi(perPageRaw)
	}

	withReceiptsRaw := values.Get("with_receipts")
	if withReceiptsRaw != "" {
		value, err := strconv.ParseBool(withReceiptsRaw)
		if err == nil {
			withReceipts = &value
		}
	}

	return &ListPaymentsRequest{
		Query:        values.Get("q"),
		UserIDs:      values["user_ids"],
		PaymentIDs:   values["member_ids"],
		WithReceipts: withReceipts,
		Page:         int64(page),
		PerPage:      int64(perPage),
		Sort:         values.Get("sort"),
	}
}

func (p *ListPaymentsRequest) Apply(query *dbx.SelectQuery) *dbx.SelectQuery {
	perPage := DefaultListCount
	if p.PerPage > 0 && p.PerPage < 500 {
		perPage = p.PerPage
	}

	page := int64(1)
	if p.Page > 0 {
		page = p.Page
	}

	sort := "issued_at desc"
	if p.Sort != "" {
		sort = p.Sort
	}

	expr := dbx.NewExp("")

	if len(p.UserIDs) > 0 {
		expr = dbx.And(dbx.In("created_by_user_id", list.ToInterfaceSlice(p.UserIDs)...))
	}

	if len(p.PaymentIDs) > 0 {
		expr = dbx.And(dbx.In("member_id", list.ToInterfaceSlice(p.PaymentIDs)...))
	}

	if p.WithReceipts != nil {
		if *p.WithReceipts {
			expr = dbx.And(dbx.NewExp("receipt_id is not null"))
		} else {
			expr = dbx.And(dbx.NewExp("receipt_id is null"))
		}
	}

	if p.Query != "" {
		query = query.InnerJoin(
			"members",
			dbx.NewExp("payments.member_id = members.id"),
		)

		searchQuery := utils.Normalize(p.Query)

		queryExpr := dbx.Or(
			dbx.Like("issued_at", searchQuery),
			dbx.Like("concat(members.first_name, ' ', members.last_name)", searchQuery),
			dbx.Like("concat(members.last_name, ' ', members.first_name)", searchQuery),
			dbx.Like("members.mobile", searchQuery),
			dbx.Like("members.phone", searchQuery),
			dbx.Like("members.email", searchQuery),
			dbx.HashExp{"members.member_no": searchQuery},
		)

		expr = dbx.And(expr, queryExpr)
	}

	return query.
		Where(expr).
		OrderBy(sort).
		Limit(perPage).
		Offset((page - 1) * perPage)
}
