package query

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/utils"
)

const DefaultListCount = int64(200)

//nolint:tagliatelle
type SearchParams struct {
	Query                   string   `json:"q"`
	MemberNo                string   `json:"member_no"`
	Name                    string   `json:"name"`
	Mobile                  string   `json:"mobile"`
	Phone                   string   `json:"phone"`
	Email                   string   `json:"email"`
	ActiveOnly              bool     `json:"active_only"`
	CompanyID               string   `json:"company_id"`
	LegacyArea              string   `json:"legacy_area"`
	AddressCityIDs          []string `json:"address_city_ids"`
	BusinessTypeIDs         []string `json:"business_type_ids"`
	WithComments            bool     `json:"with_comments"`
	ChapterID               string   `json:"chapter_id"`
	Columns                 []string `json:"columns"`
	Limit                   int64    `json:"limit"`
	WithFixedMonthlyPayment bool     `json:"with_fixed_monthly_payment"`
}

func (p *SearchParams) CompanyIDs() []string {
	return strings.Split(p.CompanyID, ",")
}

func (p *SearchParams) LegacyAreas() []string {
	areas := []string{}

	for _, area := range strings.Split(p.LegacyArea, ",") {
		area = utils.Normalize(area)
		if area == "" {
			continue
		}

		areas = append(areas, area)
	}

	return areas
}

func NewSearchParams(values url.Values) *SearchParams {
	limit := 0

	limitRaw := values.Get("limit")
	if limitRaw != "" {
		limit, _ = strconv.Atoi(limitRaw)
	}

	return &SearchParams{
		Query:                   values.Get("q"),
		MemberNo:                values.Get("member_no"),
		Name:                    values.Get("name"),
		Mobile:                  values.Get("mobile"),
		Phone:                   values.Get("phone"),
		Email:                   values.Get("email"),
		ActiveOnly:              values.Get("active_only") != "false",
		CompanyID:               values.Get("company_id"),
		LegacyArea:              values.Get("legacy_area"),
		AddressCityIDs:          values["address_city_ids"],
		BusinessTypeIDs:         values["business_type_ids"],
		WithComments:            values.Get("with_comments") == "true",
		ChapterID:               values.Get("chapter_id"),
		WithFixedMonthlyPayment: values.Get("with_fixed_monthly_payment") == "true",
		Limit:                   int64(limit),
	}
}

func (p *SearchParams) Apply(query *dbx.SelectQuery) *dbx.SelectQuery {
	limit := DefaultListCount
	if p.Limit > 0 {
		limit = p.Limit
	}

	query = query.Limit(limit)

	expr := dbx.NewExp("")

	// Searching by arbitrary terms always takes precedence.
	//nolint:nestif
	if p.Query != "" {
		searchQuery := utils.Normalize(p.Query)

		queryExpr := dbx.Or(
			dbx.Like("concat(first_name, ' ', last_name)", searchQuery),
			dbx.Like("concat(last_name, ' ', first_name)", searchQuery),
			dbx.Like("mobile", searchQuery),
			dbx.Like("phone", searchQuery),
			dbx.Like("email", searchQuery),
			dbx.HashExp{"member_no": searchQuery},
		)

		expr = dbx.And(expr, queryExpr)
	} else {
		if p.MemberNo != "" {
			expr = dbx.Or(expr, dbx.HashExp{"member_no": utils.Normalize(p.MemberNo)})
		}

		if p.Name != "" {
			expr = dbx.Or(expr, dbx.Like("full_name", utils.Normalize(p.Name)))
		}

		if p.Mobile != "" {
			expr = dbx.Or(expr, dbx.Like("mobile", utils.Normalize(p.Mobile)))
		}

		if p.Phone != "" {
			expr = dbx.Or(expr, dbx.Like("phone", utils.Normalize(p.Phone)))
		}

		if p.Email != "" {
			expr = dbx.Or(expr, dbx.Like("email", utils.Normalize(p.Email)))
		}
	}

	if p.WithComments {
		query = query.InnerJoin(
			"payments",
			dbx.NewExp("payments.member_id = members.id"),
		)

		expr = dbx.And(
			expr,
			dbx.Or(
				dbx.Not(dbx.HashExp{"members.comments": ""}),
				dbx.Not(dbx.HashExp{"payments.comments": ""}),
			),
		)
	}

	if p.ActiveOnly {
		query = query.InnerJoin(
			"subscriptions",
			dbx.NewExp("subscriptions.member_id = members.id AND subscriptions.active is true"),
		)
	}

	if p.CompanyID != "" || len(p.BusinessTypeIDs) > 0 {
		if len(p.BusinessTypeIDs) > 0 {
			query = query.InnerJoin(
				"companies",
				dbx.NewExp("members.company_id = companies.id"),
			)
		}

		companyExpr := dbx.In("members.company_id", list.ToInterfaceSlice(p.CompanyIDs())...)
		businessTypeExpr := dbx.In("companies.business_type_id", list.ToInterfaceSlice(p.BusinessTypeIDs)...)

		switch {
		case p.CompanyID != "" && len(p.BusinessTypeIDs) > 0:
			expr = dbx.And(expr, dbx.Or(companyExpr, businessTypeExpr))
		case p.CompanyID != "":
			expr = dbx.And(expr, companyExpr)
		case len(p.BusinessTypeIDs) > 0:
			expr = dbx.And(expr, businessTypeExpr)
		}
	}

	// Filtering by address city/area can be combined with the legacy field.
	// What we need to take care is group together the address related filters
	// to use OR (inclusive).
	if len(p.AddressCityIDs) > 0 || p.LegacyArea != "" {
		addressExpr := []dbx.Expression{}

		if len(p.AddressCityIDs) > 0 {
			addressExpr = append(addressExpr, dbx.In("address_city_id", list.ToInterfaceSlice(p.AddressCityIDs)...))
		}

		if p.LegacyArea != "" {
			for _, legacyArea := range p.LegacyAreas() {
				addressExpr = append(addressExpr, dbx.Like("legacy_area", legacyArea))
			}
		}

		expr = dbx.And(expr, dbx.Or(addressExpr...))
	}

	if p.WithFixedMonthlyPayment {
		expr = dbx.And(expr, dbx.NewExp("members.fixed_monthly_amount_in_euros > 0"))
	}

	return query.Distinct(true).Where(expr)
}
