package employment

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/fragoulis/setip_v2/internal/db/company"
)

func NewRandom(dao *daos.Dao, memberID string, employed bool) (*Employment, error) {
	return New(
		WithMemberID(memberID),
		WithRandomStartDate(),
		WithRandomCompany(dao),
		WithEmployed(employed),
	)
}

func WithMemberID(memberID string) Option {
	return func(e *Employment) error {
		e.MemberID = memberID

		return nil
	}
}

func WithRandomCompany(dao *daos.Dao) Option {
	return func(emp *Employment) error {
		company, err := company.Random(dao)
		if err != nil {
			return fmt.Errorf("random company: %w", err)
		}

		emp.CompanyID = company.GetId()

		return nil
	}
}

func WithEmployed(employed bool) Option {
	return func(emp *Employment) error {
		if employed {
			return nil
		}

		// Generate a date greater than the start date.
		date, err := types.ParseDateTime(gofakeit.DateRange(
			emp.StartDate.Time(),
			time.Now().AddDate(0, 0, -1),
		))
		if err != nil {
			return err
		}

		emp.EndDate = date

		if emp.StartDate.Time().Compare(emp.EndDate.Time()) >= 0 {
			return fmt.Errorf("start date must be before the end date: %s %s",
				emp.StartDate,
				emp.EndDate,
			)
		}

		return nil
	}
}

func WithRandomStartDate() Option {
	return func(emp *Employment) error {
		// Generate a date in past 10 years.
		date, err := types.ParseDateTime(gofakeit.DateRange(
			time.Now().AddDate(-10, 0, 0),
			time.Now().AddDate(-1, 0, 0),
		))
		if err != nil {
			return err
		}

		emp.StartDate = date

		return nil
	}
}
