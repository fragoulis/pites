package address

import (
	"fmt"
	"strings"

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

func FindStreetByName(dao *daos.Dao, name string, city *City) (*Street, error) {
	model := &Street{}

	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	err := StreetQuery(dao).
		Where(dbx.NewExp("name_normalized = {:name}", dbx.Params{"name": strings.ToLower(name)})).
		AndWhere(dbx.NewExp("city_id = {:city_id}", dbx.Params{"city_id": city.GetId()})).
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}

func FindCityByID(dao *daos.Dao, id string) (*City, error) {
	model := &City{}

	err := CityQuery(dao).
		Where(dbx.NewExp("id = {:id}", dbx.Params{"id": id})).
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}

func FindCityByName(dao *daos.Dao, name string) (*City, error) {
	model := &City{}

	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	err := CityQuery(dao).
		Where(dbx.NewExp("name = {:name}", dbx.Params{"name": name})).
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}
