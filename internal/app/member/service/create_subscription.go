package service

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/events"
	"github.com/fragoulis/setip_v2/internal/utils"
)

type CreateSubscriptionRequest struct {
	StartDate string `json:"start_date"`
	FeePaid   bool   `json:"fee_paid"`
}

func (r *CreateSubscriptionRequest) Validate() validation.Errors {
	errs := validation.Errors{}

	if r.StartDate == "" {
		errs["start_date"] = requiredErr
	}

	return errs
}

func CreateSubscription(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	member *models.Record,
	data *CreateSubscriptionRequest,
) error {
	// Run our custom validations and return early on failure.
	errs := data.Validate()
	if len(errs) != 0 {
		return errs
	}

	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return errors.ErrFailedToGetDao
	}

	return dao.RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		collection, err := tx.FindCollectionByNameOrId("subscriptions")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}

		newPayment := models.NewRecord(collection)

		form := forms.NewRecordUpsert(app, newPayment)
		form.SetDao(tx)

		err = form.LoadData(map[string]any{
			"member_id":          member.GetId(),
			"active":             true,
			"fee_paid":           data.FeePaid,
			"start_date":         data.StartDate,
			"created_by_user_id": utils.CurrentUserID(ctx),
		})
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		return events.WrapCreate(ctx, app, newPayment, func() error {
			return form.Submit()
		})
	})
}
