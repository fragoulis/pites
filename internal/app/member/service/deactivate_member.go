package service

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/events"
)

func DeactivateMember(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	member *models.Record,
) error {
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

		subscriptions, err := tx.FindRecordsByExpr(
			collection.GetId(),
			dbx.HashExp{"member_id": member.GetId(), "end_date": ""},
		)
		if err != nil {
			return fmt.Errorf("failed to find member: %w", err)
		}

		if len(subscriptions) == 0 {
			return nil
		}

		subscription := subscriptions[0]
		form := forms.NewRecordUpsert(app, subscription)
		form.SetDao(tx)

		err = form.LoadData(map[string]any{
			"active":   false,
			"end_date": time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		return events.WrapUpdate(ctx, app, subscription, func() (*models.Record, error) {
			err := form.Submit()
			if err != nil {
				return nil, fmt.Errorf("failed to update subscription: %w", err)
			}

			return subscription, nil
		})
	})
}
