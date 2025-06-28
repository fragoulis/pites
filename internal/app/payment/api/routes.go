package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/payment/http"
)

func RegisterRoutes(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/payments", http.List(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.GET("/payments/:id", http.Get(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.POST("/payments", http.Create(app), apis.RequireAdminOrRecordAuth("users"))
	e.Router.PATCH("/payments", http.Update(app), apis.RequireAdminOrRecordAuth("users"))
}
