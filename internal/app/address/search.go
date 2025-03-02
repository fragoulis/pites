package address

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/utils"
)

const DefaultListCount = 50

type Address struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Search(_ echo.Context, dao *daos.Dao, query string) ([]*Address, error) {
	records := []*models.Record{}
	query = strings.ToLower(utils.Normalize(query))

	sel := dao.RecordQuery("addresses")

	for _, token := range strings.Split(query, " ") {
		if token == "" {
			continue
		}

		sel = sel.AndWhere(dbx.Like("search_term", token))
	}

	err := sel.Limit(DefaultListCount).All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	addresses := []*Address{}
	for _, record := range records {
		addresses = append(addresses, NewAddressFromRecord(record))
	}

	return addresses, nil
}

func NewAddressFromRecord(rec *models.Record) *Address {
	model := &Address{
		ID: rec.GetId(),
		Name: fmt.Sprintf("%s, %s, %s",
			rec.GetString("street"),
			rec.GetString("city"),
			rec.GetString("zipcode"),
		),
	}

	return model
}

func ListCities(_ echo.Context, dao *daos.Dao, query string) ([]*City, error) {
	records := []*models.Record{}
	query = utils.Normalize(query)

	sel := dao.RecordQuery("address_cities")

	for _, token := range strings.Split(query, " ") {
		if token == "" {
			continue
		}

		sel = sel.AndWhere(dbx.Like("name", token))
	}

	err := sel.All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	cities := make([]*City, 0, len(records))
	for _, record := range records {
		cities = append(cities, &City{
			ID:   record.GetId(),
			Name: record.GetString("name"),
		})
	}

	return cities, nil
}
