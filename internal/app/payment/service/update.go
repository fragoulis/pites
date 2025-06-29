package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
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
	WithoutReceipt bool   `json:"without_receipt"`
}

type UpdatePaymentResponse CreatePaymentResponse

//nolint:gocyclo,maintidx
func Update(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	id string,
	data *UpdatePaymentRequest,
) (*UpdatePaymentResponse, error) {
	dao := app.Dao()

	paymentRecord, err := dao.FindRecordById("payments", id)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	receiptRecord, err := dao.FindRecordById("receipts", paymentRecord.GetString("receipt_id"))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db error on finding receipt: %w", err)
		}
	}

	errs := validation.Errors{}

	// If receipt does not exist and without receipt is true, ignore.
	// If receipt does not exist and without receipt is false, create it.
	// If receipt exists and without receipt is true, and only this payment refers to that receipt, delete it.
	// If receipt exists and without receipt is false, and only this payment refers to that receipt, update it.
	// If receipt exists and without receipt is true, and more payments refer to that receipt, validation error.
	// If receipt exists and without receipt is false, and more payments refer to that receipt, validation error.
	if receiptRecord != nil {
		paymentsMatchingReceipt, err := dao.FindRecordsByFilter(
			"payments",
			"receipt_id = {:receipt_id}",
			"-created",
			0,
			0,
			dbx.Params{"receipt_id": receiptRecord.GetId()},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find payments: %w", err)
		}

		if len(paymentsMatchingReceipt) > 1 {
			errs["amount"] = cannotEditPayment
		}
	}

	if data.Amount <= 0 {
		errs["amount"] = invalidAmount
	}

	if !data.WithoutReceipt && data.ReceiptNo > 0 && data.ReceiptBlockNo <= 0 {
		errs["receipt_block_no"] = requiredReceiptBlockNo
	}

	if len(errs) != 0 {
		return nil, errs
	}

	months := int(math.Floor(float64(data.Amount) / float64(subscriptionFee)))

	memberID := paymentRecord.GetString("member_id")

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

	currentMonths := paymentRecord.GetInt("months")

	// This date is calculated including the payment we are updating.
	// Subtract the payment's months from the date.
	memberHasPaidUntil = utils.EndOfMonthAhead(memberHasPaidUntil, -currentMonths)

	err = app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		//nolint:nestif
		if !data.WithoutReceipt {
			if receiptRecord == nil {
				receiptsCollection, err := tx.FindCollectionByNameOrId("receipts")
				if err != nil {
					return fmt.Errorf("failed to find receipts collection: %w", err)
				}

				receiptRecord = models.NewRecord(receiptsCollection)
			}

			form := forms.NewRecordUpsert(app, receiptRecord)
			form.SetDao(tx)

			newReceiptData := map[string]any{
				"amount_in_euros": data.Amount,
				"comments":        data.Comments,
			}
			if receiptRecord.IsNew() {
				newReceiptData["member_id"] = paymentRecord.GetString("member_id")
				newReceiptData["issued_at"] = paymentRecord.GetDateTime("issued_at")
				newReceiptData["created_by_user_id"] = utils.CurrentUserID(ctx)
			}

			if data.ReceiptNo > 0 {
				newReceiptData["receipt_no"] = data.ReceiptNo
			}

			if data.ReceiptBlockNo > 0 {
				newReceiptData["block_no"] = data.ReceiptBlockNo
			}

			err := form.LoadData(newReceiptData)
			if err != nil {
				return fmt.Errorf("failed to load receipt data: %w", err)
			}

			if receiptRecord.IsNew() {
				err = events.WrapCreate(ctx, app, receiptRecord, func() error {
					return form.Submit()
				})
			} else {
				err = events.WrapUpdate(
					ctx,
					app,
					receiptRecord,
					func() (*models.Record, error) {
						return receiptRecord, form.Submit()
					},
				)
			}

			if err != nil {
				//nolint:wrapcheck
				return err
			}
		} else if data.WithoutReceipt && receiptRecord != nil {
			// Delete the receipt record.
			err = tx.Delete(receiptRecord)
			if err != nil {
				return fmt.Errorf("failed to delete receipt: %w", err)
			}

			receiptRecord = nil
		}

		receiptID := ""
		if receiptRecord != nil {
			receiptID = receiptRecord.GetId()
		}

		form := forms.NewRecordUpsert(app, paymentRecord)
		form.SetDao(tx)

		newData := map[string]any{
			"amount_in_euros": data.Amount,
			"months":          months,
			"comments":        data.Comments,
			"receipt_id":      receiptID,
		}

		if !memberHasPaidUntil.IsZero() {
			newData["legacy_to"] = utils.EndOfMonthAhead(memberHasPaidUntil, months)
		}

		err = form.LoadData(newData)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		err = events.WrapUpdate(
			ctx,
			app,
			paymentRecord,
			func() (*models.Record, error) {
				return paymentRecord, form.Submit()
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
		Payment: model.NewFromRecordNoMember(paymentRecord, member.MemberNo, member.FullName),
		Status:  member.PaymentStatus.Formatted,
	}, nil
}
