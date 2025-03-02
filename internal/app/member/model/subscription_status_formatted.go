package model

import (
	"fmt"
	"time"

	"github.com/fragoulis/setip_v2/internal/utils"
)

type subscriptionStatusFormatted struct {
	startDate time.Time
	endDate   time.Time
	active    bool
	feePaid   bool
}

func newSubscriptionStatusFormattedFromSubscription(sub *Subscription) *subscriptionStatusFormatted {
	return newSubscriptionStatusFormatted(
		sub.StartDate,
		sub.EndDate,
		sub.Active,
		sub.FeePaid,
	)
}

func newSubscriptionStatusFormatted(
	startDate, endDate time.Time,
	active, feePaid bool,
) *subscriptionStatusFormatted {
	return &subscriptionStatusFormatted{startDate, endDate, active, feePaid}
}

func (s *subscriptionStatusFormatted) String() string {
	startDateFormatted := utils.Month(s.startDate)
	feePaidFormatted := "Πληρωμένη εγγραφή"

	if !s.feePaid {
		feePaidFormatted = "Χρωστάει εγγραφή"
	}

	if s.active {
		return fmt.Sprintf("Ενεργή/ός από %s (%s)", startDateFormatted, feePaidFormatted)
	}

	endDateFormatted := utils.Month(s.endDate)

	return fmt.Sprintf(
		"Ενεργή/ός από %s μέχρι %s",
		startDateFormatted,
		endDateFormatted,
	)
}
