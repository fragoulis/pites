package query

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/payment/model"
)

const (
	DefaultListCount = int64(50)
)

//nolint:tagliatelle
type ListPaymentsRequest struct {
	ID        string   `json:"id"`
	MemberIDs []string `json:"member_ids"`
	UserIDs   []string `json:"user_ids"`
	Limit     int64    `json:"limit"`
}

func NewListPaymentsRequest(values url.Values) *ListPaymentsRequest {
	limit := 0

	limitRaw := values.Get("limit")
	if limitRaw != "" {
		limit, _ = strconv.Atoi(limitRaw)
	}

	return &ListPaymentsRequest{
		ID:        values.Get("id"),
		UserIDs:   values["user_ids"],
		MemberIDs: values["member_ids"],
		Limit:     int64(limit),
	}
}

func (p *ListPaymentsRequest) Apply(query *dbx.SelectQuery) *dbx.SelectQuery {
	if p.ID != "" {
		return query.Where(&dbx.HashExp{"id": p.ID}).Limit(1)
	}

	limit := DefaultListCount
	if p.Limit > 0 {
		limit = p.Limit
	}

	expr := dbx.NewExp("")

	if len(p.UserIDs) > 0 {
		expr = dbx.And(dbx.In("created_by_user_id", list.ToInterfaceSlice(p.UserIDs)...))
	}

	if len(p.MemberIDs) > 0 {
		expr = dbx.And(dbx.In("member_id", list.ToInterfaceSlice(p.MemberIDs)...))
	}

	return query.Where(expr).Limit(limit)
}

func List(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	params *ListPaymentsRequest,
) ([]*model.Payment, error) {
	dao := app.Dao()

	query := params.Apply(dao.RecordQuery("payments"))

	records := []*models.Record{}

	err := query.OrderBy("issued_at DESC").All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	err = apis.EnrichRecords(
		ctx,
		dao,
		records,
		"member_id",
		"created_by_user_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to expand relations for payments: %w", err)
	}

	payments := make([]*model.Payment, 0, len(records))
	for _, record := range records {
		payments = append(payments, model.NewFromRecord(record))
	}

	return payments, nil
}
