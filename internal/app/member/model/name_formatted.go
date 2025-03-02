package model

import (
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

type nameFormatted struct {
	first  string
	last   string
	father string
}

func newNameFormattedFromRecord(rec *models.Record) *nameFormatted {
	return newNameFormatted(
		rec.GetString("first_name"),
		rec.GetString("last_name"),
		rec.GetString("father_name"),
	)
}

func newNameFormatted(first, last, father string) *nameFormatted {
	return &nameFormatted{first, last, father}
}

func (n *nameFormatted) String() string {
	return fmt.Sprintf("%s %s του %s", n.last, n.first, n.father)
}
