package subscription

import (
	"github.com/pocketbase/pocketbase"

	"github.com/fragoulis/setip_v2/internal/db/auditlog"
)

func RegisterCallbacks(app *pocketbase.PocketBase) {
	auditlog.RegisterCallbacks(app, "subscriptions")
}
