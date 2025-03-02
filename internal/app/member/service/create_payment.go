package service

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/payment/service"
)

type CreatePaymentRequest service.CreatePaymentRequest

func CreatePayment(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	member *models.Record,
	data *CreatePaymentRequest,
) error {
	request := &service.CreatePaymentRequest{
		MemberID:       member.GetId(),
		Amount:         data.Amount,
		ReceiptNo:      data.ReceiptNo,
		ReceiptBlockNo: data.ReceiptBlockNo,
	}

	_, err := service.Create(ctx, app, request)
	if err != nil {
		return fmt.Errorf("payment create: %w", err)
	}

	return nil
}
