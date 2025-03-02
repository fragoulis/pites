package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/rest"

	"github.com/fragoulis/setip_v2/internal/app/payment/query"
	"github.com/fragoulis/setip_v2/internal/app/payment/service"
)

func List(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		request := query.NewListPaymentsRequest(ctx.QueryParams())

		models, err := query.List(ctx, app, request)
		if err != nil {
			return apis.NewBadRequestError("failed to list payments", err)
		}

		if request.ID != "" {
			if len(models) > 0 {
				return ctx.JSON(http.StatusOK, models[0])
			}

			return ctx.JSON(http.StatusNotFound, nil)
		}

		return ctx.JSON(http.StatusOK, models)
	}
}

func Create(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		data := &service.CreatePaymentRequest{}

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.Create(ctx, app, data)
		if err != nil {
			return apis.NewBadRequestError("failed to create payment", err)
		}

		return ctx.JSON(http.StatusCreated, model)
	}
}

func Update(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		id := ctx.QueryParams().Get("id")
		if id == "" {
			return apis.NewBadRequestError("missing id", nil)
		}

		data := &service.UpdatePaymentRequest{}

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.Update(ctx, app, id, data)
		if err != nil {
			return apis.NewBadRequestError("failed to update payment", err)
		}

		return ctx.JSON(http.StatusOK, model)
	}
}
