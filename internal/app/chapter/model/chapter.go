package model

import (
	"github.com/pocketbase/pocketbase/models"
)

//nolint:tagliatelle
type Chapter struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	RawCityQuery string   `json:"raw_city_query"`
	CityIDs      []string `json:"city_ids"`
}

func NewFromRecord(rec *models.Record) *Chapter {
	return &Chapter{
		ID:           rec.GetId(),
		Name:         rec.GetString("name"),
		RawCityQuery: rec.GetString("raw_city_query"),
		CityIDs:      rec.GetStringSlice("city_ids"),
	}
}
