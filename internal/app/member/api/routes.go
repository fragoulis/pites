package api

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/fragoulis/setip_v2/internal/app/member/http"
)

func RegisterRoutes(srvEvnt *core.ServeEvent, app *pocketbase.PocketBase) {
	srvEvnt.Router.GET("/members", http.Search(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.GET("/members/:id", http.Get(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.POST("/members/:id/payment", http.CreatePayment(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.PATCH("/members/:id", http.UpdateMember(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.GET("/members/next", http.GetNextMemberNo(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.POST("/members", http.CreateMember(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.POST("/members/export", http.Export(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.POST("/members/import", http.Import(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.POST("/members/:id/subscriptions", http.Activate(app), apis.RequireAdminOrRecordAuth("users"))
	srvEvnt.Router.DELETE("/members/:id/subscriptions", http.Deactivate(app), apis.RequireAdminOrRecordAuth("users"))
}
