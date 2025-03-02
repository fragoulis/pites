package events

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type (
	UpdateFn func() (*models.Record, error)
	CreateFn func() error
)

func WrapUpdate(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	record *models.Record,
	updateFn UpdateFn,
) error {
	event := new(core.RecordUpdateEvent)
	event.HttpContext = ctx
	event.Collection = record.Collection()
	event.Record = record

	err := app.OnRecordBeforeUpdateRequest().Trigger(event)
	if err != nil {
		return fmt.Errorf("failed to execute before update request: %w", err)
	}

	updatedRecord, err := updateFn()
	if err != nil {
		return err
	}

	event.Record = updatedRecord

	err = app.OnRecordAfterUpdateRequest().Trigger(event)
	if err != nil {
		return fmt.Errorf("failed to execute after update request: %w", err)
	}

	return nil
}

func WrapCreate(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	record *models.Record,
	createFn CreateFn,
) error {
	event := new(core.RecordCreateEvent)
	event.HttpContext = ctx
	event.Collection = record.Collection()
	event.Record = record

	err := createFn()
	if err != nil {
		return err
	}

	err = app.OnRecordAfterCreateRequest().Trigger(event)
	if err != nil {
		return fmt.Errorf("failed to execute after create request: %w", err)
	}

	return nil
}
