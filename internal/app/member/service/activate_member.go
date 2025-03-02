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
	"github.com/fragoulis/setip_v2/internal/utils"
)

type ActivateMemberRequest struct {
	FeePaid bool `json:"fee_paid"`
}

func ActivateMember(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	member *models.Record,
	data *ActivateMemberRequest,
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

		if len(subscriptions) > 0 {
			return nil
		}

		subscription := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, subscription)
		form.SetDao(tx)

		err = form.LoadData(map[string]any{
			"member_id":          member.GetId(),
			"active":             true,
			"fee_paid":           data.FeePaid,
			"start_date":         time.Now(),
			"created_by_user_id": utils.CurrentUserID(ctx),
		})
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		return events.WrapCreate(ctx, app, subscription, func() error {
			return form.Submit()
		})
	})
}
