package address

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*Street)(nil)

type Street struct {
	models.BaseModel

	Name    string `db:"name"    json:"name"`
	Zipcode string `db:"zipcode" json:"zipcode"`
	CityID  string `db:"city_id" json:"city_id"`
}

func (m *Street) TableName() string {
	return "address_streets"
}

func StreetQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Street{})
}

func RandomStreet(dao *daos.Dao) (*Street, error) {
	model := &Street{}

	err := StreetQuery(dao).
		OrderBy("RANDOM()").
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}

func FindStreetByID(dao *daos.Dao, id string) (*Street, error) {
	model := &Street{}

	err := StreetQuery(dao).
		Where(dbx.NewExp("id = {:id}", dbx.Params{"id": id})).
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}
