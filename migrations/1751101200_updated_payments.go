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
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u6gs4b34n5sryaq")
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(`[
			"CREATE UNIQUE INDEX ` + "`" + `idx_wKn1viw` + "`" + ` ON ` + "`" + `payments` + "`" + ` (` + "`" + `legacy_guid` + "`" + `) WHERE ` + "`" + `legacy_guid` + "`" + ` != ''"
		]`), &collection.Indexes); err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("aokflly3")

		// remove
		collection.Schema.RemoveField("fjuon9z4")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u6gs4b34n5sryaq")
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(`[
			"CREATE UNIQUE INDEX ` + "`" + `idx_A0oAcq8` + "`" + ` ON ` + "`" + `payments` + "`" + ` (\n  ` + "`" + `receipt_block_no` + "`" + `,\n  ` + "`" + `receipt_no` + "`" + `\n) WHERE ` + "`" + `receipt_block_no` + "`" + ` != '' AND ` + "`" + `receipt_no` + "`" + ` != ''",
			"CREATE UNIQUE INDEX ` + "`" + `idx_wKn1viw` + "`" + ` ON ` + "`" + `payments` + "`" + ` (` + "`" + `legacy_guid` + "`" + `) WHERE ` + "`" + `legacy_guid` + "`" + ` != ''"
		]`), &collection.Indexes); err != nil {
			return err
		}

		// add
		del_receipt_block_no := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "aokflly3",
			"name": "receipt_block_no",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_receipt_block_no); err != nil {
			return err
		}
		collection.Schema.AddField(del_receipt_block_no)

		// add
		del_receipt_no := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "fjuon9z4",
			"name": "receipt_no",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_receipt_no); err != nil {
			return err
		}
		collection.Schema.AddField(del_receipt_no)

		return dao.SaveCollection(collection)
	})
}
