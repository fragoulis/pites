package service

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/issue/model"
)

func CreateIssue(dao *daos.Dao, issue *model.Issue) (bool, error) {
	collection, err := dao.FindCollectionByNameOrId("issues")
	if err != nil {
		return false, err
	}

	record, err := dao.FindFirstRecordByFilter(
		collection.GetId(),
		"issue_type_id = {:issue_type_id} && relation_name = {:relation_name} && relation_id = {:relation_id}",
		dbx.Params{"issue_type_id": issue.IssueTypeID},
		dbx.Params{"relation_name": issue.RelationName},
		dbx.Params{"relation_id": issue.RelationID},
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	// Create only if it does not already exist.
	created := false

	if record == nil {
		created = true

		record = models.NewRecord(collection)
		record.Set("issue_type_id", issue.IssueTypeID)
		record.Set("relation_name", issue.RelationName)
		record.Set("relation_id", issue.RelationID)
	}

	// Always update the tangible columns.
	record.Set("importance", issue.Importance)
	record.Set("resolved_at", "")

	err = dao.SaveRecord(record)
	if err != nil {
		return false, err
	}

	return created, nil
}
