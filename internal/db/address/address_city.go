package address

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*City)(nil)

type City struct {
	models.BaseModel

	Name string `db:"name" json:"name"`
}

func (m *City) TableName() string {
	return "address_cities"
}

func CityQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&City{})
}

func RandomCity(dao *daos.Dao) (*City, error) {
	model := &City{}

	err := CityQuery(dao).
		OrderBy("RANDOM()").
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}
