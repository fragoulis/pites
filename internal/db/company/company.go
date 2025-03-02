package company

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*Company)(nil)

type Company struct {
	models.BaseModel

	Name            string `db:"name"`
	Branch          string `db:"branch"`
	Email           string `db:"email"`
	Phone           string `db:"phone"`
	Website         string `db:"website"`
	AddressCityID   string `db:"address_city_id"`
	AddressStreetID string `db:"address_street_id"`
	AddressStreetNo string `db:"address_street_no"`
	LegacyAddress   string `db:"legacy_address"`
	LegacyGUID      string `db:"legacy_guid"`
	SearchTerms     string `db:"search_terms"`
}

func (m *Company) TableName() string {
	return "companies"
}

type Option func(*Company) error

func New(options ...Option) (*Company, error) {
	model := &Company{}

	for _, option := range options {
		err := option(model)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func Query(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Company{})
}

func Random(dao *daos.Dao) (*Company, error) {
	model := &Company{}

	err := Query(dao).
		OrderBy("RANDOM()").
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}
