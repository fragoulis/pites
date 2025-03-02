package subscription

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/tools/types"
)

func NewRandom(memberID string, paid, active bool) (*Subscription, error) {
	return New(
		WithMemberID(memberID),
		WithRandomStartDate(),
		WithPaidFee(paid),
		WithActive(active),
	)
}

func WithMemberID(memberID string) Option {
	return func(s *Subscription) error {
		s.MemberID = memberID

		return nil
	}
}

func WithPaidFee(paid bool) Option {
	return func(s *Subscription) error {
		s.FeePaid = paid

		return nil
	}
}

func WithActive(active bool) Option {
	return func(sub *Subscription) error {
		sub.Active = active

		if active {
			return nil
		}

		// Generate a date greater than the start date.
		date, err := types.ParseDateTime(gofakeit.DateRange(
			sub.StartDate.Time(),
			time.Now().AddDate(0, 0, -1),
		))
		if err != nil {
			return err
		}

		sub.EndDate = date

		if sub.StartDate.Time().Compare(sub.EndDate.Time()) >= 0 {
			return fmt.Errorf("start date must be before the end date: %s %s",
				sub.StartDate,
				sub.EndDate,
			)
		}

		return nil
	}
}

func WithRandomStartDate() Option {
	return func(sub *Subscription) error {
		// Generate a date in past 10 years.
		date, err := types.ParseDateTime(gofakeit.DateRange(
			time.Now().AddDate(-10, 0, 0),
			time.Now().AddDate(-1, 0, 0),
		))
		if err != nil {
			return err
		}

		sub.StartDate = date

		return nil
	}
}
