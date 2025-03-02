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

		collection, err := dao.FindCollectionByNameOrId("companies")
		if err != nil {
			return err
		}

		// add
		newBusinessTypeID := &schema.SchemaField{}

		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "lnjkrsl3",
			"name": "business_type_id",
			"type": "relation",
			"required": false,
			"presentable": true,
			"unique": false,
			"options": {
				"collectionId": "cl4z4jhnr6y571t",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), newBusinessTypeID); err != nil {
			return err
		}

		collection.Schema.AddField(newBusinessTypeID)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("companies")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("lnjkrsl3")

		return dao.SaveCollection(collection)
	})
}
