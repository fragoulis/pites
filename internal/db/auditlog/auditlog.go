package auditlog

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func insert(
	app core.App,
	ctx echo.Context,
	record *models.Record,
	event string,
	before *models.Record,
) error {
	if record.Collection().Name == "auditlogs" {
		// ignore auditlog record changes
		return nil
	}

	// This is a dirty way to wrap auditlogs in transitions
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		// Fallback, because this is also run from the admin panel
		dao = app.Dao()
	}

	auditlogCollection, err := dao.FindCollectionByNameOrId("auditlogs")
	if err != nil {
		return fmt.Errorf("failed to find audit collection: %w", err)
	}

	var adminID, authRecordID string

	// get the authenticated admin
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	if admin != nil {
		adminID = admin.Id
	}

	// or get the authenticated user record
	authRecord, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	if authRecord != nil {
		authRecordID = authRecord.Id
	}

	auditlog := models.NewRecord(auditlogCollection)
	auditlog.Set("table", record.Collection().Name)
	auditlog.Set("event", event)
	auditlog.Set("record_id", record.GetId())
	auditlog.Set("user_id", authRecordID)
	auditlog.Set("admin_id", adminID)
	auditlog.Set("before", before)
	auditlog.Set("after", record)

	err = dao.SaveRecord(auditlog)
	if err != nil {
		return fmt.Errorf("failed to save auditlog: %w", err)
	}

	return nil
}

func RegisterCallbacks(app *pocketbase.PocketBase, collection string) {
	app.OnRecordAfterCreateRequest(collection).
		Add(func(e *core.RecordCreateEvent) error {
			return insert(
				app,
				e.HttpContext,
				e.Record,
				"create",
				nil,
			)
		})

	var before *models.Record

	app.OnRecordBeforeUpdateRequest(collection).Add(
		func(e *core.RecordUpdateEvent) error {
			fmt.Println("HOST: ", e.HttpContext.Request().Host)
			before = e.Record.CleanCopy()

			return nil
		},
	)

	app.OnRecordAfterUpdateRequest(collection).
		Add(func(e *core.RecordUpdateEvent) error {
			return insert(
				app,
				e.HttpContext,
				e.Record,
				"update",
				before,
			)
		})
}
