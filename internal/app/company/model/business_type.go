package model

import (
	"github.com/pocketbase/pocketbase/models"
)

type BusinessType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewBusinessTypeFromRecord(rec *models.Record) *BusinessType {
	if rec == nil {
		return nil
	}

	return &BusinessType{
		ID:   rec.GetId(),
		Name: rec.GetString("name"),
	}
}
