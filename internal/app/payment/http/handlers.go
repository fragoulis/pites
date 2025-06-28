package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/rest"

	"github.com/fragoulis/setip_v2/internal/app/payment/model"
	"github.com/fragoulis/setip_v2/internal/app/payment/query"
	"github.com/fragoulis/setip_v2/internal/app/payment/service"
)

type ListResponse struct {
	Records []*model.Payment `json:"records"`
	Total   int              `json:"total"`
}

func List(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		request := query.NewListPaymentsRequest(ctx.QueryParams())

		ctx.Set("dao", app.Dao())

		models, err := query.List(ctx, request)
		if err != nil {
			return apis.NewBadRequestError("failed to list payments", err)
		}

		count, err := query.Count(ctx, request)
		if err != nil {
			return apis.NewBadRequestError("failed to count payments", err)
		}

		return ctx.JSON(http.StatusOK, &ListResponse{
			Records: models,
			Total:   count,
		})
	}
}

func Get(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		model, err := query.FindByID(ctx, id)
		if err != nil {
			return apis.NewNotFoundError("failed to find payment", err)
		}

		return ctx.JSON(http.StatusOK, model)
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
