package service

import (
	"errors"
	"fmt"
	"math"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	memberModel "github.com/fragoulis/setip_v2/internal/app/member/model"
	"github.com/fragoulis/setip_v2/internal/app/member/query"
	"github.com/fragoulis/setip_v2/internal/app/payment/model"
	"github.com/fragoulis/setip_v2/internal/events"
	"github.com/fragoulis/setip_v2/internal/utils"
)

type UpdatePaymentRequest struct {
	Amount         int    `json:"amount"`
	ReceiptBlockNo int    `json:"receipt_block_no"`
	ReceiptNo      int    `json:"receipt_no"`
	Comments       string `json:"comments"`
}

type UpdatePaymentResponse CreatePaymentResponse

func Update(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	id string,
	data *UpdatePaymentRequest,
) (*UpdatePaymentResponse, error) {
	dao := app.Dao()

	record, err := dao.FindRecordById("payments", id)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	errs := validation.Errors{}

	if data.Amount <= 0 {
		errs["amount"] = invalidAmount
	}

	if data.ReceiptNo > 0 && data.ReceiptBlockNo <= 0 {
		errs["receipt_block_no"] = requiredReceiptBlockNo
	}

	if len(errs) != 0 {
		return nil, errs
	}

	months := int(math.Floor(float64(data.Amount) / float64(subscriptionFee)))

	memberID := record.GetString("member_id")

	member, err := query.FindByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to find member: %w", err)
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

	currentMonths := record.GetInt("months")

	// This date is calculated including the payment we are updating.
	// Subtract the payment's months from the date.
	memberHasPaidUntil = utils.EndOfMonthAhead(memberHasPaidUntil, -currentMonths)

	err = app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		form := forms.NewRecordUpsert(app, record)
		form.SetDao(tx)

		newData := map[string]any{
			"amount_in_euros": data.Amount,
			"months":          months,
			"comments":        data.Comments,
		}

		if !memberHasPaidUntil.IsZero() {
			newData["legacy_to"] = utils.EndOfMonthAhead(memberHasPaidUntil, months)
		}

		if data.ReceiptNo > 0 {
			newData["receipt_no"] = data.ReceiptNo
		}

		if data.ReceiptBlockNo > 0 {
			newData["receipt_block_no"] = data.ReceiptBlockNo
		}

		err := form.LoadData(newData)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		err = events.WrapUpdate(
			ctx,
			app,
			record,
			func() (*models.Record, error) {
				err := form.Submit()
				if err != nil {
					return record, err
				}

				return record, nil
			},
		)
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		member, err = query.FindByID(ctx, memberID)
		if err != nil {
			return fmt.Errorf("failed to find member: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &UpdatePaymentResponse{
		Payment: model.NewFromRecordNoMember(record, member.MemberNo, member.FullName),
		Status:  member.PaymentStatus.Formatted,
	}, nil
}
