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

const (
	subscriptionFee = 2
	RegistrationFee = 4
)

//nolint:gochecknoglobals
var memberNotFoundErr = validation.NewError("validation_required", "Δεν υπάρχει μέλος με αυτό το ID")

//nolint:gochecknoglobals
var insufficientAmount = validation.NewError("validation_required", "Το ποσό δεν επαρκεί αφού περιέχει εγγραφή")

//nolint:gochecknoglobals
var requiredReceiptBlockNo = validation.NewError(
	"validation_required",
	"Το πεδίο πρέπει να συμπληρωθεί αφού υπάρχει αριθμός απόδειξης.",
)

//nolint:gochecknoglobals
var cannotEditPayment = validation.NewError(
	"validation_required",
	"Η είσπραξη ή/και η απόδειξη δε μπορούν να τροποποιηθούν γιατί η απόδειξη είναι κοινή για παραπάνω από μία εισπράξεις.",
)

//
//nolint:gochecknoglobals
var invalidAmount = validation.NewError(
	"validation_required",
	"Είτε το ποσό είτε οι μήνες πρέπει να είναι μεγαλύτερο του μηδέν (0)",
)

//nolint:gochecknoglobals
var invalidMonths = validation.NewError("validation_required", "Οι μήνες δε μπορούν να υπερβαίνουν τους 24")

type CreatePaymentRequest struct {
	MemberID                string `json:"member_id"`
	Amount                  int    `json:"amount"`
	Months                  int    `json:"months"`
	ContainsRegistrationFee bool   `json:"contains_registration_fee"`
	ReceiptBlockNo          int    `json:"receipt_block_no"`
	ReceiptNo               int    `json:"receipt_no"`
	IssuedAt                string `json:"issued_at"`
	Comments                string `json:"comments"`
	WithoutReceipt          bool   `json:"without_receipt"`
}

type CreatePaymentResponse struct {
	*model.Payment
	Status string `json:"status"`
}

//nolint:funlen,gocognit,cyclop
func Create(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	data *CreatePaymentRequest,
) (*CreatePaymentResponse, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, internalErrors.ErrFailedToGetDao
	}

	errs := validation.Errors{}

	// Make sure that the member id is valid.
	_, err := dao.FindRecordById("members", data.MemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errs["member_id"] = memberNotFoundErr
		}
	}

	if data.Months <= 0 {
		if data.Amount <= 0 {
			errs["amount"] = invalidAmount
		} else if data.ContainsRegistrationFee && data.Amount < RegistrationFee {
			errs["amount"] = insufficientAmount
		}

		if !data.WithoutReceipt && data.ReceiptNo > 0 && data.ReceiptBlockNo <= 0 {
			errs["receipt_block_no"] = requiredReceiptBlockNo
		}
	}

	//nolint:mnd
	if data.Months > 24 {
		errs["months"] = invalidMonths
	}

	if len(errs) != 0 {
		return nil, errs
	}

	// We control the payment's value by either the amount or the months passed.
	// Months take precedence.
	// If months are positive, amount is set to zero.
	// If not and amount is positive, months is set to amount/fee.
	if data.Months > 0 {
		data.Amount = 0
	} else {
		data.Months = int(math.Floor(float64(data.Amount) / float64(subscriptionFee)))
	}

	paymentsCollection, err := app.Dao().FindCollectionByNameOrId("payments")
	if err != nil {
		return nil, fmt.Errorf("failed to find payments collection: %w", err)
	}

	newPayment := models.NewRecord(paymentsCollection)

	receiptsCollection, err := app.Dao().FindCollectionByNameOrId("receipts")
	if err != nil {
		return nil, fmt.Errorf("failed to find receipts collection: %w", err)
	}

	member, err := query.FindByID(ctx, data.MemberID)
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

	// Add time to this date.
	parsedIssuedAt, err := time.Parse(time.DateOnly, data.IssuedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue date: %w", err)
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

	if data.ContainsRegistrationFee {
		data.Months -= 2
	}

	err = dao.RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		newReceipt := models.NewRecord(receiptsCollection)

		if !data.WithoutReceipt {
			form := forms.NewRecordUpsert(app, newReceipt)
			form.SetDao(tx)

			newReceiptData := map[string]any{
				"member_id":          data.MemberID,
				"amount_in_euros":    data.Amount,
				"issued_at":          issuedAt,
				"comments":           data.Comments,
				"created_by_user_id": utils.CurrentUserID(ctx),
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

			err = events.WrapCreate(ctx, app, newPayment, func() error {
				return form.Submit()
			})
			if err != nil {
				//nolint:wrapcheck
				return err
			}
		}

		form := forms.NewRecordUpsert(app, newPayment)
		form.SetDao(tx)

		newPaymentData := map[string]any{
			"member_id":                 data.MemberID,
			"amount_in_euros":           data.Amount,
			"months":                    data.Months,
			"contains_registration_fee": data.ContainsRegistrationFee,
			"issued_at":                 issuedAt,
			"comments":                  data.Comments,
			"created_by_user_id":        utils.CurrentUserID(ctx),
			"receipt_id":                newReceipt.GetId(),
		}

		if !memberHasPaidUntil.IsZero() {
			newPaymentData["legacy_to"] = utils.EndOfMonthAhead(memberHasPaidUntil, data.Months)
		}

		err := form.LoadData(newPaymentData)
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

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		Payment: model.NewFromRecordNoMember(newPayment, member.MemberNo, member.FullName),
		Status:  member.PaymentStatus.Formatted,
	}, nil
}
