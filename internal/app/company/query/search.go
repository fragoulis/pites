package query

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/utils"
)

const DefaultListCount = 50

func Search(ctx echo.Context, dao *daos.Dao, request *ListCompaniesRequest) ([]*model.Company, error) {
	records := []*models.Record{}

	query := dao.RecordQuery("companies")

	if request != nil && request.Query != "" {
		query = query.Where(dbx.Like("name", utils.Normalize(request.Query)))
	}

	if request != nil && len(request.BusinessTypeIDs) > 0 {
		query = query.AndWhere(dbx.In("business_type_id", list.ToInterfaceSlice(request.BusinessTypeIDs)...))
	}

	limit := DefaultListCount
	if request != nil && request.Limit > 0 {
		limit = request.Limit
	}

	err := query.
		Limit(int64(limit)).
		OrderBy("name ASC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Expand the records relations (aka load associations).
	err = apis.EnrichRecords(
		ctx,
		dao,
		records,
		"parent_id",
		"address_city_id",
		"address_street_id",
		"business_type_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to expand relations for companies: %w", err)
	}

	companies := []*model.Company{}

	for _, record := range records {
		parentID := record.GetString("parent_id")
		if parentID != "" {
			parent := record.ExpandedOne("parent_id")

			// If parent id is present, the record is a branch,
			// in which case we swap the name/branch columns.
			record.Set("branch", record.GetString("name"))
			record.Set("name", parent.GetString("name"))
		}

		// Map the pocketbase company to our custom model.
		companies = append(companies, model.NewFromRecord(record))
	}

	return companies, nil
}

func Count(_ echo.Context, dao *daos.Dao, request *ListCompaniesRequest) (int, error) {
	var count int

	query := dao.RecordQuery("companies").
		Where(dbx.Like("name", utils.Normalize(request.Query)))

	if len(request.BusinessTypeIDs) > 0 {
		query = query.AndWhere(dbx.In("business_type_id", list.ToInterfaceSlice(request.BusinessTypeIDs)...))
	}

	err := query.Select("count(*)").
		Limit(1).Row(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return count, nil
}

func FindByID(ctx echo.Context, id string) (*model.Company, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	record, err := dao.FindRecordById("companies", id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Expand the records relations (aka load associations).
	err = apis.EnrichRecord(
		ctx,
		dao,
		record,
		"address_city_id",
		"address_street_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to expand relations for companies: %w", err)
	}

	return model.NewFromRecord(
		record,
	), nil
}

func ListBusinessTypes(_ echo.Context, dao *daos.Dao) ([]*model.BusinessType, error) {
	records := []*models.Record{}

	err := dao.RecordQuery("company_business_types").OrderBy("name ASC").All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	businessTypes := make([]*model.BusinessType, 0, len(records))
	for _, record := range records {
		businessTypes = append(businessTypes, model.NewBusinessTypeFromRecord(record))
	}

	return businessTypes, nil
}
