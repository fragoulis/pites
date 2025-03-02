package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/address/http"
)

func RegisterRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/address", http.Search(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.GET("/address_cities", http.ListCities(app), apis.RequireAdminOrRecordAuth("users"))
}
