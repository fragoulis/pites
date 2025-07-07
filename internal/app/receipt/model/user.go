package model

import "github.com/pocketbase/pocketbase/models"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func newUserFromRecord(rec *models.Record) *User {
	if rec == nil {
		return nil
	}

	return &User{
		ID:       rec.GetId(),
		Username: rec.GetString("username"),
	}
}
