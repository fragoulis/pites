package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "cl4z4jhnr6y571t",
			"created": "2024-10-26 08:17:31.291Z",
			"updated": "2024-10-26 08:17:31.291Z",
			"name": "company_business_types",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "tfmspjs3",
					"name": "name",
					"type": "text",
					"required": true,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_u9Mejj4` + "`" + ` ON ` + "`" + `company_business_types` + "`" + ` (` + "`" + `name` + "`" + `)"
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		dao := daos.New(db)

		if err := dao.SaveCollection(collection); err != nil {
			return err
		}

		types := map[string]string{
			// reducted
		}

		for name, id := range types {
			record := models.NewRecord(collection)
			record.Set("id", id)
			record.Set("name", name)

			if err := dao.SaveRecord(record); err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("cl4z4jhnr6y571t")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
