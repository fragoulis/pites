package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/rest"

	"github.com/fragoulis/setip_v2/internal/app/receipt/service"
)

func Create(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		data := &service.CreateReceiptRequest{}

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.Create(ctx, app, data)
		if err != nil {
			return apis.NewBadRequestError("failed to create receipt", err)
		}

		return ctx.JSON(http.StatusCreated, model)
	}
}

func CreateBatch(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		data := &service.CreateBatchReceiptRequest{}

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.CreateBatch(ctx, app, data)
		if err != nil {
			return apis.NewBadRequestError("failed to create receipts", err)
		}

		return ctx.JSON(http.StatusCreated, model)
	}
}
