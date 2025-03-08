package service

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/member/model"
	"github.com/fragoulis/setip_v2/internal/events"
)

func CreateMember(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	dao *daos.Dao,
	data *CreateMemberRequest,
) (*model.Member, error) {
	// Run our custom validations and return early on failure.
	errs := data.Validate()
	if len(errs) != 0 {
		return nil, errs
	}

	var newRec *models.Record

	err := dao.RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		collection, err := tx.FindCollectionByNameOrId("members")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}

		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		form.SetDao(tx)

		formData, err := data.ToUpdateRequest().ToFormData(tx)
		if err != nil {
			return fmt.Errorf("failed to form data: %w", err)
		}

		// We set this explicitly here. We do not want it be part of
		// the updateable payload.
		formData["member_no"] = data.MemberNo

		err = form.LoadData(formData)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		// Create member
		err = events.WrapCreate(
			ctx,
			app,
			record,
			func() error {
				return form.Submit()
			},
		)
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		// Create subscription
		err = CreateSubscription(ctx, app, record, &CreateSubscriptionRequest{
			FeePaid:   data.FeePaid,
			StartDate: data.StartDate,
		})
		if err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}

		// Fetch newly created record
		newRec, err = tx.FindRecordById("members", record.GetId())
		if err != nil {
			return fmt.Errorf("failed to find new record: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.NewFromRecord(newRec, nil, nil, nil, nil), nil
}
