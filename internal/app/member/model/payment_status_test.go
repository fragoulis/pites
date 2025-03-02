package model_test

import (
	"testing"
	"time"

	"github.com/fragoulis/setip_v2/internal/app/member/model"
	"github.com/fragoulis/setip_v2/internal/utils"
)

func TestPaymentStatus(t *testing.T) {
	utils.SetTimeNow(time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name              string
		paidUntil         time.Time
		expectedFormatted string
		expectedOK        bool
	}{
		{
			name:              "member is ok",
			paidUntil:         time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "ΟΚ μέχρι και Νοέμβριος του 2024.",
			expectedOK:        true,
		},
		{
			name:              "member is ok",
			paidUntil:         time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "ΟΚ μέχρι και Νοέμβριος του 2024.",
			expectedOK:        true,
		},
		{
			name:              "member is ok",
			paidUntil:         time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "ΟΚ μέχρι και Ιανουάριος του 2025.",
			expectedOK:        true,
		},
		{
			name:              "member is not ok",
			paidUntil:         time.Date(2024, 10, 18, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "Χρωστάει 2 € (1 μήνες). ΟΚ μέχρι και Οκτώβριος του 2024.",
			expectedOK:        false,
		},
		{
			name:              "member is not ok",
			paidUntil:         time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "Χρωστάει 2 € (1 μήνες). ΟΚ μέχρι και Οκτώβριος του 2024.",
			expectedOK:        false,
		},
		{
			name:              "member is not ok",
			paidUntil:         time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "Χρωστάει 10 € (5 μήνες). ΟΚ μέχρι και Ιούνιος του 2024.",
			expectedOK:        false,
		},
		{
			name:              "member is not ok",
			paidUntil:         time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			expectedFormatted: "Χρωστάει 34 € (17 μήνες). ΟΚ μέχρι και Ιούνιος του 2023.",
			expectedOK:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := model.NewPaymentStatusWithActiveSubscription(test.paidUntil)
			if test.expectedFormatted != status.Formatted {
				t.Errorf("Expected %q but got %q", test.expectedFormatted, status.Formatted)
			}

			if test.expectedOK != status.OK {
				t.Errorf("Expected %t but got %t", test.expectedOK, status.OK)
			}
		})
	}
}
