package address

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/db/address"
)

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Street struct {
	address.Street
}

func Build(rec *models.Record) string {
	street := rec.ExpandedOne("address_street_id")
	city := rec.ExpandedOne("address_city_id")

	if street == nil || city == nil {
		return ""
	}

	return fmt.Sprintf("%s %s, %s, %s",
		street.GetString("name"),
		rec.GetString("address_street_no"),
		city.GetString("name"),
		street.GetString("zipcode"),
	)
}

func FindStreetsByID(dao *daos.Dao, ids []string) ([]*Street, error) {
	streets := []*Street{}

	err := address.StreetQuery(dao).
		Where(dbx.In("id", list.ToInterfaceSlice(ids)...)).
		OrderBy("name ASC").
		All(&streets)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return streets, nil
}

func FindCitiesByID(dao *daos.Dao, ids []string) ([]*City, error) {
	cities := []*City{}

	err := address.CityQuery(dao).
		Where(dbx.In("id", list.ToInterfaceSlice(ids)...)).
		OrderBy("name ASC").
		All(&cities)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return cities, nil
}
