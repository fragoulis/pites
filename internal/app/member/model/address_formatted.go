package model

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

type addressFormatted struct {
	street   string
	streetNo string
	city     string
	zipcode  string
}

func newAddressFormattedFromRecord(rec *models.Record) string {
	street := rec.ExpandedOne("address_street_id")
	city := rec.ExpandedOne("address_city_id")

	if street == nil && city == nil {
		legacyFields := []string{}

		for _, field := range []string{
			rec.GetString("legacy_address"),
			rec.GetString("legacy_area"),
			rec.GetString("legacy_city"),
			rec.GetString("legacy_post_code"),
		} {
			if field != "" {
				legacyFields = append(legacyFields, field)
			}
		}

		return strings.Join(legacyFields, ", ")
	}

	if street == nil {
		return city.GetString("name")
	}

	return newAddressFormatted(
		street.GetString("name"),
		rec.GetString("address_street_no"),
		city.GetString("name"),
		street.GetString("zipcode"),
	).String()
}

func newAddressStreetNameFromRecord(rec *models.Record) string {
	city := rec.ExpandedOne("address_street_id")

	if city == nil {
		return ""
	}

	return city.GetString("name")
}

func newAddressCityNameFromRecord(rec *models.Record) string {
	city := rec.ExpandedOne("address_city_id")

	if city == nil {
		return ""
	}

	return city.GetString("name")
}

func newAddressFormatted(street, streetNo, city, zipcode string) *addressFormatted {
	return &addressFormatted{street, streetNo, city, zipcode}
}

func (n *addressFormatted) String() string {
	if n.street == "" {
		return ""
	}

	return fmt.Sprintf("%s %s, %s, %s",
		n.street,
		n.streetNo,
		n.city,
		n.zipcode,
	)
}
