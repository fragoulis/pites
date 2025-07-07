package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/receipt/http"
)

func RegisterRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.POST("/receipts", http.Create(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.POST("/receipts/batch", http.CreateBatch(app), apis.RequireAdminOrRecordAuth("users"))
}
