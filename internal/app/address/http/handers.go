package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/fragoulis/setip_v2/internal/app/address"
)

func Search(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		query := ctx.QueryParam("q")

		addresses, err := address.Search(ctx, app.Dao(), query)
		if err != nil {
			return apis.NewBadRequestError("failed to search addresses", err)
		}

		return ctx.JSON(http.StatusOK, addresses)
	}
}

func ListCities(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		query := ctx.QueryParam("q")

		cities, err := address.ListCities(ctx, app.Dao(), query)
		if err != nil {
			return apis.NewBadRequestError("failed to list cities", err)
		}

		return ctx.JSON(http.StatusOK, cities)
	}
}
