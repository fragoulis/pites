package employment

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	companyModel "github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/utils"
)

const (
	daysInMonth = 30
	hoursInDay  = 24
)

type Employment struct {
	ID                 string                `json:"id"`
	MemberID           string                `json:"member_id"`
	CompanyID          string                `json:"company_id"`
	StartDate          time.Time             `json:"start_date"`
	EndDate            time.Time             `json:"end_date"`
	Company            *companyModel.Company `json:"company"`
	Branch             *companyModel.Company `json:"branch"`
	StartDateFormatted string                `json:"start_date_formatted"`
	EndDateFormatted   string                `json:"end_date_formatted"`
}

func (e *Employment) StillEmployed() bool {
	return e.EndDate.IsZero()
}

func (e *Employment) Months() int {
	return e.MonthsFrom(e.StartDate)
}

func (e *Employment) MonthsFrom(from time.Time) int {
	var timePassed time.Duration

	if e.StillEmployed() {
		timePassed = time.Since(from)
	} else {
		timePassed = e.EndDate.Sub(from)
	}

	return int(timePassed.Hours() / hoursInDay / daysInMonth)
}

func NewFromRecord(rec *models.Record) *Employment {
	model := &Employment{
		ID:                 rec.GetId(),
		MemberID:           rec.GetString("member_id"),
		CompanyID:          rec.GetString("company_id"),
		StartDate:          rec.GetDateTime("start_date").Time(),
		EndDate:            rec.GetDateTime("end_date").Time(),
		StartDateFormatted: utils.Month(rec.GetDateTime("start_date").Time()),
		EndDateFormatted:   utils.Month(rec.GetDateTime("end_date").Time()),
	}

	companyRec := rec.ExpandedOne("company_id")
	if companyRec != nil {
		model.Company = companyModel.NewFromRecord(companyRec)
	}

	branchRec := rec.ExpandedOne("branch_id")
	if branchRec != nil {
		model.Branch = companyModel.NewFromRecord(branchRec)
	}

	return model
}

func FindByMemberID(ctx echo.Context, memberIDs []string) ([]*Employment, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("employments").
		Where(dbx.In("member_id", list.ToInterfaceSlice(memberIDs)...)).
		OrderBy("end_date ASC", "start_date DESC").
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
			"company_id.business_type_id",
			"branch_id",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to expand relations for members: %w", err)
		}
	} else {
		errs := dao.ExpandRecords(records, []string{"company_id"}, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to expand relations for members: %v", errs)
		}
	}

	models := []*Employment{}
	for _, record := range records {
		models = append(models, NewFromRecord(record))
	}

	return models, nil
}
