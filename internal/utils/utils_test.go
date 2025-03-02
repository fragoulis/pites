package utils_test

import (
	"testing"
	"time"

	"github.com/fragoulis/setip_v2/internal/utils"
)

//nolint:dupl
func TestBeginningOfMonth(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{
			input:    time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-01-01",
		},
		{
			input:    time.Date(2000, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-02-01",
		},
		{
			input:    time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-03-01",
		},
		{
			input:    time.Date(2000, 4, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-04-01",
		},
		{
			input:    time.Date(2000, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-05-01",
		},
		{
			input:    time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-06-01",
		},
		{
			input:    time.Date(2000, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-07-01",
		},
		{
			input:    time.Date(2000, 8, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-08-01",
		},
		{
			input:    time.Date(2000, 9, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-09-01",
		},
		{
			input:    time.Date(2000, 10, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-10-01",
		},
		{
			input:    time.Date(2000, 11, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-11-01",
		},
		{
			input:    time.Date(2000, 12, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-12-01",
		},
	}

	for _, test := range tests {
		actual := utils.BeginningOfMonth(test.input).Format(time.DateOnly)

		if actual != test.expected {
			t.Errorf("Expected %q but got %q", test.expected, actual)
		}
	}
}

//nolint:dupl
func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{
			input:    time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-01-31",
		},
		{
			input:    time.Date(2000, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-02-29",
		},
		{
			input:    time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-03-31",
		},
		{
			input:    time.Date(2000, 4, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-04-30",
		},
		{
			input:    time.Date(2000, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-05-31",
		},
		{
			input:    time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-06-30",
		},
		{
			input:    time.Date(2000, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-07-31",
		},
		{
			input:    time.Date(2000, 8, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-08-31",
		},
		{
			input:    time.Date(2000, 9, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-09-30",
		},
		{
			input:    time.Date(2000, 10, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-10-31",
		},
		{
			input:    time.Date(2000, 11, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-11-30",
		},
		{
			input:    time.Date(2000, 12, 15, 0, 0, 0, 0, time.UTC),
			expected: "2000-12-31",
		},
	}

	for _, test := range tests {
		actual := utils.EndOfMonth(test.input).Format(time.DateOnly)

		if actual != test.expected {
			t.Errorf("Expected %q but got %q", test.expected, actual)
		}
	}
}

func TestMonthsSince(t *testing.T) {
	now := time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		input    time.Time
		expected int
	}{
		{
			input:    time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: -19,
		},
		{
			input:    time.Date(2021, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: -11,
		},
		{
			input:    time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: -12,
		},
		{
			input:    time.Date(2021, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: -13,
		},

		{
			input:    time.Date(2018, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 24,
		},
		{
			input:    time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: 25,
		},

		{
			input:    time.Date(2018, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: 23,
		},
		{
			input:    time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: 4,
		},
		{
			input:    time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
		{
			input:    time.Date(2020, 5, 14, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			input:    time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			input:    time.Date(2020, 5, 16, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			input:    time.Date(2020, 6, 14, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			input:    time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			input:    time.Date(2020, 6, 16, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			input:    time.Date(2020, 7, 14, 0, 0, 0, 0, time.UTC),
			expected: -1,
		},
		{
			input:    time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: -1,
		},
		{
			input:    time.Date(2020, 7, 16, 0, 0, 0, 0, time.UTC),
			expected: -1,
		},
	}

	for _, test := range tests {
		actual, _ := utils.MonthsSince(test.input, now)

		if actual != test.expected {
			t.Errorf("Expected %d but got %d", test.expected, actual)
		}
	}
}

func TestEndOfMonthAhead(t *testing.T) {
	tests := []struct {
		input    time.Time
		months   int
		expected string
	}{
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-02-28",
		},
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   2,
			expected: "2022-03-31",
		},
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   3,
			expected: "2022-04-30",
		},
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: "2021-12-31",
		},
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   -2,
			expected: "2021-11-30",
		},
		{
			input:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			months:   -3,
			expected: "2021-10-31",
		},
		{
			input:    time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-02-28",
		},
		{
			input:    time.Date(2022, 1, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-02-28",
		},
		{
			input:    time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2023-01-31",
		},
		{
			input:    time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			months:   2,
			expected: "2023-02-28",
		},
		{
			input:    time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-03-31",
		},
		{
			input:    time.Date(2022, 10, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-11-30",
		},
		{
			input:    time.Date(2022, 10, 31, 0, 0, 0, 0, time.UTC),
			months:   2,
			expected: "2022-12-31",
		},
		{
			input:    time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: "2022-12-31",
		},
		{
			input:    time.Date(2022, 10, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: "2022-09-30",
		},
		{
			input:    time.Date(2022, 3, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: "2022-02-28",
		},
	}

	for _, test := range tests {
		actual := utils.EndOfMonthAhead(test.input, test.months).Format(time.DateOnly)

		if actual != test.expected {
			t.Errorf(
				"%q plus %d months should be %q but is %q",
				test.input.Format(time.DateOnly),
				test.months,
				test.expected,
				actual,
			)
		}
	}
}
