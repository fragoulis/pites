package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/fragoulis/setip_v2/internal/app/chapter/query"
)

func List(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		chapters, err := query.Search(ctx, app.Dao())
		if err != nil {
			return apis.NewBadRequestError("failed to list chapters", err)
		}

		return ctx.JSON(http.StatusOK, chapters)
	}
}
