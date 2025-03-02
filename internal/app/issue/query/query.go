package query

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/app/issue/model"
)

func FindUnresolvedByRelationID(ctx echo.Context, collection string, ids []string) ([]*model.Issue, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("issues").
		Where(dbx.HashExp{"relation_name": collection}).
		Where(dbx.In("relation_id", list.ToInterfaceSlice(ids)...)).
		Where(dbx.HashExp{"resolved_at": ""}).
		OrderBy("created ASC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Expand the records relations (aka load associations).

	if ctx.Request() != nil {
		err = apis.EnrichRecords(
			ctx,
			dao,
			records,
			"issue_type_id",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to expand relations for issues: %w", err)
		}
	} else {
		errs := dao.ExpandRecords(records, []string{"issue_type_id"}, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to expand relations for issues: %v", errs)
		}
	}

	models := []*model.Issue{}
	for _, record := range records {
		models = append(models, model.NewFromRecord(record))
	}

	return models, nil
}

func FindIssueTypesByKey(dao *daos.Dao, keys ...string) ([]*model.IssueType, error) {
	records := []*models.Record{}

	query := dao.RecordQuery("issue_types")

	if len(keys) > 0 {
		query = query.Where(dbx.In("key", list.ToInterfaceSlice(keys)...))
	}

	err := query.All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	models := []*model.IssueType{}
	for _, record := range records {
		models = append(models, model.NewIssueTypeFromRecord(record))
	}

	return models, nil
}

func MapIssueTypesByKey(dao *daos.Dao, keys ...string) (model.IssueTypesByKey, error) {
	issueTypes, err := FindIssueTypesByKey(dao, keys...)
	if err != nil {
		return nil, err
	}

	issueTypeByKey := model.IssueTypesByKey{}

	for _, issueType := range issueTypes {
		issueTypeByKey[issueType.Key] = issueType
	}

	return issueTypeByKey, nil
}

func FindIssueTypeIDsByKey(dao *daos.Dao, keys ...string) ([]string, error) {
	issueTypes, err := FindIssueTypesByKey(dao, keys...)
	if err != nil {
		return nil, err
	}

	issueTypeIDs := make([]string, len(issueTypes))

	for i, issueType := range issueTypes {
		issueTypeIDs[i] = issueType.ID
	}

	return issueTypeIDs, nil
}
