package service

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/events"
	"github.com/fragoulis/setip_v2/internal/utils"
)

func CreateInternal(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	dao *daos.Dao,
	payment *models.Record,
	data *CreateReceiptRequest,
) (string, error) {

	collection, err := dao.FindCollectionByNameOrId("receipts")
	if err != nil {
		return "", fmt.Errorf("failed to find collection: %w", err)
	}

	newReceipt := models.NewRecord(collection)

	// Add time to this date.
	parsedIssuedAt, err := time.Parse(time.DateOnly, data.IssuedAt)
	if err != nil {
		return "", fmt.Errorf("failed to parse issue date: %w", err)
	}

	// Get the current time
	currentTime := time.Now()

	// Update the parsed time with the current hour, minute, and second
	issuedAt := time.Date(
		parsedIssuedAt.Year(),
		parsedIssuedAt.Month(),
		parsedIssuedAt.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
		currentTime.Nanosecond(),
		parsedIssuedAt.Location(),
	)

	err = dao.RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		form := forms.NewRecordUpsert(app, newReceipt)
		form.SetDao(tx)

		newData := map[string]any{
			"member_id":          payment.GetString("member_id"),
			"amount_in_euros":    payment.Get("amount_in_euros"),
			"receipt_no":         data.ReceiptNo,
			"block_no":           data.BlockNo,
			"issued_at":          issuedAt,
			"comments":           data.Comments,
			"created_by_user_id": utils.CurrentUserID(ctx),
		}

		err := form.LoadData(newData)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		//nolint:wrapcheck
		err = events.WrapCreate(ctx, app, newReceipt, func() error {
			return form.Submit()
		})
		if err != nil {
			if _, ok := err.(validation.Errors); ok {
				return err
			}

			return fmt.Errorf("failed to create receipt: %w", err)
		}

		return events.WrapUpdate(ctx, app, payment, func() (*models.Record, error) {
			payment.Set("receipt_id", newReceipt.GetId())

			err := tx.SaveRecord(payment)
			if err != nil {
				return nil, fmt.Errorf("failed to update payment with receipt: %w", err)
			}

			return payment, nil
		})
	})
	if err != nil {
		return "", err
	}

	return newReceipt.GetId(), nil
}
