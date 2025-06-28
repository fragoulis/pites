package query

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/app/payment/model"
)

func List(
	ctx echo.Context,
	params *ListPaymentsRequest,
) ([]*model.Payment, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	query := params.Apply(dao.RecordQuery("payments"))

	records := []*models.Record{}

	err := query.All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	err = apis.EnrichRecords(
		ctx,
		dao,
		records,
		"member_id",
		"receipt_id",
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

func Count(ctx echo.Context, params *ListPaymentsRequest) (int, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return 0, errors.ErrFailedToGetDao
	}

	query := params.Apply(dao.RecordQuery("payments"))

	var count int

	err := query.
		Select("count(distinct payments.id)").
		Limit(1).Row(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return count, nil
}

func FindByID(ctx echo.Context, id string) (*model.Payment, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	record, err := dao.FindRecordById("payments", id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	err = apis.EnrichRecord(
		ctx,
		dao,
		record,
		"member_id",
		"receipt_id",
		"created_by_user_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to expand relations for payment: %w", err)
	}

	return model.NewFromRecord(record), nil
}
