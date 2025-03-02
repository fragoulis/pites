package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/goodsign/monday"

	"github.com/fragoulis/setip_v2/internal/utils"
)

const (
	subscriptionFee = 2
	daysInMonth     = 30
	hoursInDay      = 24
)

//nolint:gochecknoglobals
var cutoffDate = time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)

type PaymentStatus struct {
	TotalPaid                 int    `json:"total_paid"`
	Owed                      int    `json:"owed"`
	OK                        bool   `json:"ok"`
	Formatted                 string `json:"formatted"`
	CutoffDateFormatted       string `json:"cutoff_date_formatted"`
	RegisteredAtFormatted     string `json:"registered_at_formatted"`
	LastPaymentUntilFormatted string `json:"last_payment_until_formatted"`
	IsPaymentDisabled         bool   `json:"is_payment_disabled"`
}

func newPaymentStatusFromMember(member *Member) *PaymentStatus {
	paidUntil, err := member.MemberHasPaidUntil()
	if err != nil {
		if errors.Is(err, ErrUnableToDeterminePaymentStatus) {
			return newPaymentStatusNoActiveSubscription()
		}

		return &PaymentStatus{
			Formatted: fmt.Sprintf("unknown error: %s", err),
		}
	}

	status := NewPaymentStatusWithActiveSubscription(paidUntil)
	if status == nil {
		return &PaymentStatus{
			Formatted: "weird error: start time is zero",
		}
	}

	status.CutoffDateFormatted = monday.Format(cutoffDate, "January 2006", monday.LocaleElGR)

	activeSubscriptionStartedAt := member.ActiveSubscriptionStartedAt()
	if activeSubscriptionStartedAt.IsZero() {
		status.RegisteredAtFormatted = "-"
	} else {
		status.RegisteredAtFormatted = monday.Format(
			member.ActiveSubscriptionStartedAt(),
			"January 2006",
			monday.LocaleElGR,
		)
	}

	lastPaymentUntil := member.LastPaymentUtil()
	if lastPaymentUntil.IsZero() {
		status.LastPaymentUntilFormatted = "-"
	} else {
		status.LastPaymentUntilFormatted = monday.Format(member.LastPaymentUtil(), "January 2006", monday.LocaleElGR)
	}

	return status
}

func newPaymentStatusNoActiveSubscription() *PaymentStatus {
	return &PaymentStatus{
		Formatted:         "Ανενεργό μέλος. Χρειάζεται επανεγγραφή.",
		IsPaymentDisabled: true,
	}
}

func NewPaymentStatusWithActiveSubscription(
	paidUntil time.Time,
) *PaymentStatus {
	monthsOwed, err := utils.MonthsSinceNow(paidUntil)
	if err != nil {
		return nil
	}

	paidUntilFormatted := monday.Format(
		paidUntil,
		"ΟΚ μέχρι και January του 2006.",
		monday.LocaleElGR,
	)

	if monthsOwed <= 0 {
		return &PaymentStatus{
			Formatted: paidUntilFormatted,
			OK:        true,
		}
	}

	return &PaymentStatus{
		Owed: monthsOwed * subscriptionFee,
		Formatted: fmt.Sprintf("Χρωστάει %d € (%d μήνες). %s",
			monthsOwed*subscriptionFee,
			monthsOwed,
			paidUntilFormatted,
		),
	}
}
