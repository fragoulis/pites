package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/company/http"
)

func RegisterRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/companies", http.Search(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.GET("/companies/:id", http.Get(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.PUT("/companies/:id", http.Update(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.POST("/companies", http.Create(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.GET("/companies/business_types", http.ListBusinessTypes(app), apis.RequireAdminOrRecordAuth("users"))
}
