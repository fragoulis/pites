package service

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/events"
)

func UpdateMember(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	member *models.Record,
	data *UpdateMemberRequest,
) error {
	// Run our custom validations and return early on failure.
	errs := data.Validate()
	if len(errs) != 0 {
		return errs
	}

	return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		form := forms.NewRecordUpsert(app, member)
		form.SetDao(tx)

		formData, err := data.ToFormData(tx)
		if err != nil {
			return fmt.Errorf("failed to form data: %w", err)
		}

		err = form.LoadData(formData)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		return events.WrapUpdate(
			ctx,
			app,
			member,
			func() (*models.Record, error) {
				err := form.Submit()
				if err != nil {
					return member, err
				}

				// Update record in order for the diff to appear correctly in the audit logs
				for key, value := range formData {
					member.Set(key, value)
				}

				return member, nil
			},
		)
	})
}
