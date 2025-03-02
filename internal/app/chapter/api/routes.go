package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/chapter/http"
)

func RegisterRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/chapters", http.List(app), apis.RequireAdminOrRecordAuth("users"))
}
