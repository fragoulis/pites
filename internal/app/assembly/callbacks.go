package assembly

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterCallbacks(app *pocketbase.PocketBase) {
	app.OnRecordsListRequest("assemblies").Add(func(evt *core.RecordsListEvent) error {
		items := []*Assembly{}
		for _, rec := range evt.Records {
			items = append(items, NewFromRecord(rec))
		}

		evt.Result.Items = items

		return nil
	})
}
