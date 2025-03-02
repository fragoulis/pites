package model

import (
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

type Company struct {
	ID               string        `json:"id"`
	ParentID         string        `json:"parent_id"`
	Name             string        `json:"name"`
	Branch           string        `json:"branch"`
	Email            string        `json:"email"`
	Phone            string        `json:"phone"`
	Website          string        `json:"website"`
	AddressFormatted string        `json:"address_formatted"`
	NameFormatted    string        `json:"name_formatted"`
	AddressStreetID  string        `json:"address_street_id"`
	AddressStreetNo  string        `json:"address_street_no"`
	BusinessTypeID   string        `json:"business_type_id"`
	BusinessType     *BusinessType `json:"business_type"`
	Parent           *Company      `json:"parent"`
}

func NewFromRecord(rec *models.Record) *Company {
	if rec == nil {
		return nil
	}

	company := &Company{
		ID:              rec.GetId(),
		ParentID:        rec.GetString("parent_id"),
		Name:            rec.GetString("name"),
		Branch:          rec.GetString("branch"),
		Email:           rec.GetString("email"),
		Phone:           rec.GetString("phone"),
		Website:         rec.GetString("website"),
		AddressStreetID: rec.GetString("address_street_id"),
		AddressStreetNo: rec.GetString("address_street_no"),
		BusinessTypeID:  rec.GetString("business_type_id"),
	}

	company.formatAddress(rec)
	company.formatName()
	company.BusinessType = NewBusinessTypeFromRecord(rec.ExpandedOne("business_type_id"))
	company.Parent = NewFromRecord(rec.ExpandedOne("parent_id"))

	return company
}

func (c *Company) formatAddress(rec *models.Record) {
	street := rec.ExpandedOne("address_street_id")
	city := rec.ExpandedOne("address_city_id")

	if street == nil || city == nil {
		return
	}

	c.AddressFormatted = fmt.Sprintf("%s %s, %s, %s",
		street.GetString("name"),
		rec.GetString("address_street_no"),
		city.GetString("name"),
		street.GetString("zipcode"),
	)
}

func (c *Company) formatName() {
	if c.Branch == "" {
		c.NameFormatted = c.Name

		return
	}

	c.NameFormatted = fmt.Sprintf("%s | %s", c.Name, c.Branch)
}
