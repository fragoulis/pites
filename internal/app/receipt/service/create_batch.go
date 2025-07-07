package service

import (
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	internalErrors "github.com/fragoulis/setip_v2/internal/app/errors"
)

//nolint:gochecknoglobals
var paymentHasReceiptErr = validation.NewError("validation_required", "Η είσπραξη έχει ήδη απόξειξη")

type CreateBatchReceiptRequest struct {
	PaymentIDs []string `json:"payment_ids"`
	BlockNos   []int    `json:"block_nos"`
	ReceiptNos []int    `json:"receipt_nos"`
	IssuedAt   string   `json:"issued_at"`
	Comments   string   `json:"comments"`
}

type CreateBatchReceiptResponse struct {
	IDs []string `json:"ids"`
}

//nolint:funlen,gocognit,cyclop
func CreateBatch(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	data *CreateBatchReceiptRequest,
) (*CreateBatchReceiptResponse, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, internalErrors.ErrFailedToGetDao
	}

	payments := []*models.Record{}

	errs := validation.Errors{}
	paymentIDErrs := validation.Errors{}
	receiptNoErrs := validation.Errors{}
	blockNoErrs := validation.Errors{}

	for i := range data.PaymentIDs {
		index := fmt.Sprintf("%d", i)

		if data.PaymentIDs[i] == "" {
			paymentIDErrs[index] = requiredErr
		}

		// Load the payment.
		payment, err := dao.FindRecordById("payments", data.PaymentIDs[i])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				paymentIDErrs[index] = paymentNotFoundErr
			}

			return nil, fmt.Errorf("failed to find payment: %w", err)
		}

		if payment.GetString("receipt_id") != "" {
			paymentIDErrs[index] = paymentHasReceiptErr
		}

		payments = append(payments, payment)

		if data.ReceiptNos[i] == 0 {
			receiptNoErrs[index] = requiredErr
		}

		if data.BlockNos[i] == 0 {
			blockNoErrs[index] = requiredErr
		}
	}

	if data.IssuedAt == "" {
		errs["issued_at"] = requiredErr
	}

	if len(errs) != 0 || len(paymentIDErrs) != 0 || len(receiptNoErrs) != 0 || len(blockNoErrs) != 0 {
		errs["payment_ids"] = paymentIDErrs
		errs["receipt_nos"] = receiptNoErrs
		errs["block_nos"] = blockNoErrs

		return nil, errs
	}

	createdIDs := []string{}

	err := dao.RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		for i, payment := range payments {
			request := &CreateReceiptRequest{
				PaymentID: data.PaymentIDs[i],
				BlockNo:   data.BlockNos[i],
				ReceiptNo: data.ReceiptNos[i],
				IssuedAt:  data.IssuedAt,
				Comments:  data.Comments,
			}

			id, err := CreateInternal(ctx, app, tx, payment, request)
			if err != nil {
				if validationErrors, ok := err.(validation.Errors); ok {
					index := fmt.Sprintf("%d", i)

					for field, fieldError := range validationErrors {
						switch field {
						case "block_no":
							blockNoErrs[index] = fieldError
						case "receipt_no":
							receiptNoErrs[index] = fieldError
						}
					}

					continue
				}

				return fmt.Errorf("failed to create receipt internal: %w", err)
			}

			createdIDs = append(createdIDs, id)
		}

		if len(errs) != 0 || len(receiptNoErrs) != 0 || len(blockNoErrs) != 0 {
			errs["payment_ids"] = paymentIDErrs
			errs["receipt_nos"] = receiptNoErrs
			errs["block_nos"] = blockNoErrs

			return errs
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreateBatchReceiptResponse{
		IDs: createdIDs,
	}, nil
}
