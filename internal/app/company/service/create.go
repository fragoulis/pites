package service

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/db/address"
)

var ErrStreetIDCast = errors.New("failed to cast street id to string")

func Create(ctx echo.Context, app *pocketbase.PocketBase, data map[string]any) (*model.Company, error) {
	// Fill the city id
	streetID, ok := data["address_street_id"].(string)
	if !ok {
		return nil, ErrStreetIDCast
	}

	if streetID != "" {
		street, err := address.FindStreetByID(app.Dao(), streetID)
		if err != nil {
			return nil, fmt.Errorf("failed to find street: %w", err)
		}

		data["address_city_id"] = street.CityID
	}

	collection, err := app.Dao().FindCollectionByNameOrId("companies")
	if err != nil {
		return nil, fmt.Errorf("failed to find collection: %w", err)
	}

	var newRec *models.Record

	err = app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		form.SetDao(tx)

		err = form.LoadData(data)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		app.Logger().Debug("create company", "data", data)

		event := new(core.RecordCreateEvent)
		event.HttpContext = ctx
		event.Collection = collection
		event.Record = record

		err = form.Submit()
		if err != nil {
			return err
		}

		err = app.OnRecordAfterCreateRequest().Trigger(event)
		if err != nil {
			return fmt.Errorf("failed to execute after create request: %w", err)
		}

		// Fetch newly created record
		newRec, err = tx.FindRecordById("companies", record.GetId())
		if err != nil {
			return fmt.Errorf("failed to find new record: %w", err)
		}

		err = apis.EnrichRecord(
			ctx,
			tx,
			newRec,
			"address_city_id",
			"address_street_id",
		)
		if err != nil {
			return fmt.Errorf("failed to enrich record: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.NewFromRecord(newRec), nil
}
