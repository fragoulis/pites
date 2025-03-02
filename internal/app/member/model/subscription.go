package model

import (
	"time"

	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/utils"
)

type Subscription struct {
	ID                 string    `json:"id"`
	MemberID           string    `json:"member_id"`
	Active             bool      `json:"active"`
	FeePaid            bool      `json:"fee_paid"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	StartDateFormatted string    `json:"start_date_formatted"`
	EndDateFormatted   string    `json:"end_date_formatted"`
	StatusFormatted    string    `json:"status_formatted"`
	Months             int       `json:"months"`
}

func NewSubscriptionFromRecord(rec *models.Record) *Subscription {
	sub := &Subscription{
		ID:        rec.GetId(),
		MemberID:  rec.GetString("member_id"),
		Active:    rec.GetBool("active"),
		FeePaid:   rec.GetBool("fee_paid"),
		StartDate: rec.GetDateTime("start_date").Time(),
		EndDate:   rec.GetDateTime("end_date").Time(),
	}

	endDate := sub.EndDate
	if endDate.IsZero() {
		endDate = time.Now()
	}

	sub.StartDateFormatted = utils.Month(sub.StartDate)
	sub.EndDateFormatted = utils.Month(endDate)
	sub.Months, _ = utils.MonthsSince(sub.StartDate, endDate)
	sub.StatusFormatted = newSubscriptionStatusFormattedFromSubscription(sub).String()

	return sub
}
