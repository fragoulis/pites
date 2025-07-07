package service

import (
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"

	internalErrors "github.com/fragoulis/setip_v2/internal/app/errors"
)

//nolint:gochecknoglobals
var paymentNotFoundErr = validation.NewError("validation_required", "Δεν υπάρχει πληρωμή με αυτό το ID")

//nolint:gochecknoglobals
var requiredErr = validation.NewError(
	"validation_required",
	"Το πεδίο πρέπει να συμπληρωθεί.",
)

type CreateReceiptRequest struct {
	PaymentID string `json:"payment_id"`
	BlockNo   int    `json:"block_no"`
	ReceiptNo int    `json:"receipt_no"`
	IssuedAt  string `json:"issued_at"`
	Comments  string `json:"comments"`
}

type CreateReceiptResponse struct {
	ID string `json:"id"`
}

//nolint:funlen,gocognit,cyclop
func Create(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	data *CreateReceiptRequest,
) (*CreateReceiptResponse, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, internalErrors.ErrFailedToGetDao
	}

	errs := validation.Errors{}

	// Load the payment.
	payment, err := dao.FindRecordById("payments", data.PaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errs["payment_id"] = paymentNotFoundErr
		}
	}

	if data.ReceiptNo > 0 {
		errs["receipt_no"] = requiredErr
	}

	if data.BlockNo > 0 {
		errs["block_no"] = requiredErr
	}

	if data.IssuedAt != "" {
		errs["issued_at"] = requiredErr
	}

	if len(errs) != 0 {
		return nil, errs
	}

	id, err := CreateInternal(ctx, app, dao, payment, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipt internal: %w", err)
	}

	return &CreateReceiptResponse{
		ID: id,
	}, nil
}
