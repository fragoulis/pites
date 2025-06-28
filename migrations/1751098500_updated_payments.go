package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("u6gs4b34n5sryaq")
		if err != nil {
			return err
		}

		// add
		new_receipt_id := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "iaosck6e",
			"name": "receipt_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "pny9yt9xff7iuf2",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_receipt_id); err != nil {
			return err
		}
		collection.Schema.AddField(new_receipt_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("u6gs4b34n5sryaq")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("iaosck6e")

		return dao.SaveCollection(collection)
	})
}
