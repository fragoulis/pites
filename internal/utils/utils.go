package utils

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/goodsign/monday"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))

	// Remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	value, _, _ = transform.String(t, value)

	return value
}

func Equal(a, b string) bool {
	return Normalize(a) == Normalize(b)
}

func Year(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return monday.Format(
		value,
		"2006",
		monday.LocaleElGR,
	)
}

func Month(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return monday.Format(
		value,
		"January 2006",
		monday.LocaleElGR,
	)
}

func Day(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return monday.Format(
		value,
		"02 January 2006",
		monday.LocaleElGR,
	)
}

func ForInput(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return monday.Format(
		value,
		"2006-01-02",
		monday.LocaleElGR,
	)
}

func MonthsSince(start, end time.Time) (int, error) {
	if start.IsZero() || end.IsZero() {
		return 0, errors.New("start and end must be valid")
	}

	years := end.Year() - start.Year()
	months := end.Month() - start.Month()

	return years*12 + int(months), nil
}

func MonthsSinceNow(start time.Time) (int, error) {
	return MonthsSince(start, Now())
}

func CurrentUserID(ctx echo.Context) string {
	// get the authenticated admin
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	if admin != nil {
		return admin.Id
	}

	// or get the authenticated user record
	authRecord, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	if authRecord != nil {
		return authRecord.Id
	}

	return ""
}

func BeginningOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func EndOfMonth(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month()+1,
		0,
		0,
		0,
		0,
		0,
		date.Location(),
	)
}

func EndOfMonthAhead(date time.Time, months int) time.Time {
	return time.Date(
		date.Year(),
		date.Month()+time.Month(months)+1,
		0,
		0,
		0,
		0,
		0,
		date.Location(),
	)
}

//nolint:gochecknoglobals
var timeNow time.Time

func SetTimeNow(tf time.Time) {
	timeNow = tf
}

// Now is used in places that have tests that require
// time manipulation.
func Now() time.Time {
	if timeNow.IsZero() {
		return time.Now()
	}

	return timeNow
}
