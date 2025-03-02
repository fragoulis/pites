package query

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/chapter/model"
)

func Search(_ echo.Context, dao *daos.Dao) ([]*model.Chapter, error) {
	records := []*models.Record{}

	err := dao.RecordQuery("chapters").
		OrderBy("name ASC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	chapters := make([]*model.Chapter, 0, len(records))

	for _, record := range records {
		chapters = append(chapters, model.NewFromRecord(record))
	}

	return chapters, nil
}

func FindByID(dao *daos.Dao, id string) (*model.Chapter, bool, error) {
	records := []*models.Record{}

	err := dao.RecordQuery("chapters").Where(dbx.HashExp{"id": id}).All(&records)
	if err != nil {
		return nil, false, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(records) == 0 {
		return nil, false, nil
	}

	return model.NewFromRecord(records[0]), true, nil
}
