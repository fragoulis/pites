package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/rest"

	"github.com/fragoulis/setip_v2/internal/app/company/model"
	"github.com/fragoulis/setip_v2/internal/app/company/query"
	"github.com/fragoulis/setip_v2/internal/app/company/service"
)

type SearchResponse struct {
	Records []*model.Company `json:"records"`
	Total   int              `json:"total"`
}

func Search(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		request := query.NewListCompaniesRequestFromQueryParams(ctx.QueryParams())

		count, err := query.Count(ctx, app.Dao(), request)
		if err != nil {
			return apis.NewBadRequestError("failed to count members", err)
		}

		companies, err := query.Search(ctx, app.Dao(), request)
		if err != nil {
			return apis.NewBadRequestError("failed to search companies", err)
		}

		return ctx.JSON(http.StatusOK, &SearchResponse{
			Records: companies,
			Total:   count,
		})
	}
}

func Get(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		company, err := query.FindByID(ctx, id)
		if err != nil {
			return apis.NewNotFoundError("failed to find company", err)
		}

		return ctx.JSON(http.StatusOK, company)
	}
}

func Update(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		record, err := app.Dao().FindRecordById("companies", id)
		if err != nil {
			return apis.NewBadRequestError("failed to find company", err)
		}

		data := map[string]any{}

		err = rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.Update(ctx, app, record, data)
		if err != nil {
			return apis.NewBadRequestError("failed to update company", err)
		}

		return ctx.JSON(http.StatusOK, model)
	}
}

func Create(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		data := map[string]any{}

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		model, err := service.Create(ctx, app, app.Dao(), data)
		if err != nil {
			return apis.NewBadRequestError("failed to create company", err)
		}

		return ctx.JSON(http.StatusCreated, model)
	}
}

func ListBusinessTypes(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		businessTypes, err := query.ListBusinessTypes(ctx, app.Dao())
		if err != nil {
			return apis.NewBadRequestError("failed to list business types", err)
		}

		return ctx.JSON(http.StatusOK, businessTypes)
	}
}
