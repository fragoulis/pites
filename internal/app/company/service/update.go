package service

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/db/address"
	"github.com/fragoulis/setip_v2/internal/events"
)

func Update(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	record *models.Record,
	data map[string]any,
) (*model.Company, error) {
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

	var updatedRec *models.Record

	err := app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		form := forms.NewRecordUpsert(app, record)
		form.SetDao(tx)

		err := form.LoadData(data)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		err = events.WrapUpdate(
			ctx,
			app,
			record,
			func() (*models.Record, error) {
				err := form.Submit()
				if err != nil {
					return nil, err
				}

				// Fetch newly created record
				updatedRec, err = tx.FindRecordById("companies", record.GetId())
				if err != nil {
					return nil, fmt.Errorf("failed to find new record: %w", err)
				}

				return updatedRec, nil
			})
		if err != nil {
			return fmt.Errorf("copmany update: %w", err)
		}

		err = apis.EnrichRecord(
			ctx,
			tx,
			updatedRec,
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

	return model.NewFromRecord(updatedRec), nil
}
