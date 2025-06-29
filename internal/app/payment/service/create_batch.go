package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	internalErrors "github.com/fragoulis/setip_v2/internal/app/errors"
	memberModel "github.com/fragoulis/setip_v2/internal/app/member/model"
	"github.com/fragoulis/setip_v2/internal/app/member/query"
	"github.com/fragoulis/setip_v2/internal/app/payment/model"
	"github.com/fragoulis/setip_v2/internal/events"
	"github.com/fragoulis/setip_v2/internal/utils"
)

type CreateBatchPaymentRequest struct {
	//nolint:tagliatelle
	MemberIDs []string `json:"member_ids"`
	Amount    int      `json:"amount"`
	IssuedAt  string   `json:"issued_at"`
	Comments  string   `json:"comments"`
}

type CreateBatchPaymentResponse struct {
	Payments []*model.Payment
}

//nolint:funlen,gocognit,cyclop
func CreateBatch(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	data *CreateBatchPaymentRequest,
) (*CreateBatchPaymentResponse, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, internalErrors.ErrFailedToGetDao
	}

	errs := validation.Errors{}

	for _, memberID := range data.MemberIDs {
		// Make sure that members exist.
		_, err := dao.FindRecordById("members", memberID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errs["amount"] = memberNotFoundErr

				break
			}
		}
	}

	if data.Amount <= 0 {
		errs["amount"] = invalidAmount
	}

	if len(errs) != 0 {
		return nil, errs
	}

	paymentsCollection, err := app.Dao().FindCollectionByNameOrId("payments")
	if err != nil {
		return nil, fmt.Errorf("failed to find payments collection: %w", err)
	}

	newPayments := []*model.Payment{}

	for _, memberID := range data.MemberIDs {
		err = dao.RunInTransaction(func(tx *daos.Dao) error {
			ctx.Set("dao", tx)

			member, err := query.FindByID(ctx, memberID)
			if err != nil {
				return fmt.Errorf("failed to find member: %w", err)
			}

			// Temporary: set the legacy_to field which is needed to calculate
			// until which date a member has paid for.
			memberHasPaidUntil, err := member.MemberHasPaidUntil()
			if err != nil {
				if errors.Is(err, memberModel.ErrUnableToDeterminePaymentStatus) {
					data.Comments += `
SYSTEM: Payment created without a legacy_to because the member had no active subscription.
`
				} else {
					data.Comments += fmt.Sprintf(`
Unknown error: %s`, err)
				}
			}

			newPayment := models.NewRecord(paymentsCollection)

			// Add time to this date.
			parsedIssuedAt, err := time.Parse(time.DateOnly, data.IssuedAt)
			if err != nil {
				return fmt.Errorf("failed to parse issue date: %w", err)
			}

			// Get the current time
			currentTime := time.Now()

			// Update the parsed time with the current hour, minute, and second
			issuedAt := time.Date(
				parsedIssuedAt.Year(),
				parsedIssuedAt.Month(),
				parsedIssuedAt.Day(),
				currentTime.Hour(),
				currentTime.Minute(),
				currentTime.Second(),
				currentTime.Nanosecond(),
				parsedIssuedAt.Location(),
			)

			form := forms.NewRecordUpsert(app, newPayment)
			form.SetDao(tx)

			months := int(math.Floor(float64(data.Amount) / float64(subscriptionFee)))

			newPaymentData := map[string]any{
				"member_id":          memberID,
				"amount_in_euros":    data.Amount,
				"months":             months,
				"issued_at":          issuedAt,
				"comments":           data.Comments,
				"created_by_user_id": utils.CurrentUserID(ctx),
			}

			if !memberHasPaidUntil.IsZero() {
				newPaymentData["legacy_to"] = utils.EndOfMonthAhead(memberHasPaidUntil, months)
			}

			err = form.LoadData(newPaymentData)
			if err != nil {
				return fmt.Errorf("failed to load payment data: %w", err)
			}

			err = events.WrapCreate(ctx, app, newPayment, func() error {
				return form.Submit()
			})
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			newPayments = append(newPayments, model.NewFromRecordNoMember(newPayment, member.MemberNo, member.FullName))

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &CreateBatchPaymentResponse{
		Payments: newPayments,
	}, nil
}
