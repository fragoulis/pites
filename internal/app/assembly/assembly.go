package assembly

import (
	"time"

	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/utils"
)

type Assembly struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	DateFormatted string    `json:"date_formatted"`
	Comments      string    `json:"comments"`
	Active        bool      `json:"active"`
}

func NewFromRecord(rec *models.Record) *Assembly {
	return &Assembly{
		ID:            rec.GetId(),
		Date:          rec.GetDateTime("date").Time(),
		DateFormatted: utils.Day(rec.GetDateTime("date").Time()),
		Comments:      rec.GetString("comments"),
		Active:        rec.GetBool("active"),
	}
}
