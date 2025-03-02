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

		collection, err := dao.FindCollectionByNameOrId("members")
		if err != nil {
			return err
		}

		// add
		newCompanyID := &schema.SchemaField{}

		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "f1zz6o8l",
			"name": "company_id",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "ofjm22yoatvxtz3",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), newCompanyID); err != nil {
			return err
		}

		collection.Schema.AddField(newCompanyID)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("members")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("f1zz6o8l")

		return dao.SaveCollection(collection)
	})
}
