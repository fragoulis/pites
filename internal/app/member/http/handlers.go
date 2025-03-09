package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/rest"

	chapterQuery "github.com/fragoulis/setip_v2/internal/app/chapter/query"
	"github.com/fragoulis/setip_v2/internal/app/member/model"
	"github.com/fragoulis/setip_v2/internal/app/member/query"
	"github.com/fragoulis/setip_v2/internal/app/member/service"
)

type SearchResponse struct {
	Records []*model.Member `json:"records"`
	Total   int             `json:"total"`
}

func Search(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		searchParams := query.NewSearchParams(ctx.QueryParams())

		searchParams, err := prepareMembersearch(app.Dao(), searchParams)
		if err != nil {
			return apis.NewBadRequestError("failed to prepare query", err)
		}

		count, err := query.Count(ctx, searchParams)
		if err != nil {
			return apis.NewBadRequestError("failed to count members", err)
		}

		members, err := query.Search(ctx, searchParams)
		if err != nil {
			return apis.NewBadRequestError("failed to search members", err)
		}

		return ctx.JSON(http.StatusOK, &SearchResponse{
			Records: members,
			Total:   count,
		})
	}
}

func Get(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		member, err := query.FindByID(ctx, id)
		if err != nil {
			return apis.NewNotFoundError("failed to find member", err)
		}

		return ctx.JSON(http.StatusOK, member)
	}
}

func CreatePayment(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		record, err := app.Dao().FindRecordById("members", id)
		if err != nil {
			return apis.NewBadRequestError("failed to find member", err)
		}

		var data service.CreatePaymentRequest

		err = rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		err = service.CreatePayment(ctx, app, record, &data)
		if err != nil {
			return apis.NewBadRequestError("failed to create payment", err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}
}

func UpdateMember(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		record, err := app.Dao().FindRecordById("members", id)
		if err != nil {
			return apis.NewBadRequestError("failed to find member", err)
		}

		var data service.UpdateMemberRequest

		err = rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		_, err = service.UpdateMember(ctx, app, app.Dao(), record, &data)
		if err != nil {
			return apis.NewBadRequestError("failed to update details", err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}
}

func GetNextMemberNo(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		memberNo, err := query.NextMemberNo(ctx)
		if err != nil {
			return apis.NewBadRequestError("failed to get next member no", err)
		}

		return ctx.JSON(http.StatusOK, memberNo)
	}
}

func CreateMember(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		var data service.CreateMemberRequest

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		_, err = service.CreateMember(ctx, app, app.Dao(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to create member", err)
		}

		return ctx.JSON(http.StatusCreated, nil)
	}
}

func Export(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		var searchParamsPre query.SearchParams

		err := rest.CopyJsonBody(ctx.Request(), &searchParamsPre)
		if err != nil {
			return apis.NewBadRequestError("failed to prepare query", err)
		}

		searchParams, err := prepareMembersearch(app.Dao(), &searchParamsPre)
		if err != nil {
			return apis.NewBadRequestError("failed to prepare query", err)
		}

		searchParams.Limit = 10000

		members, err := query.Search(ctx, searchParams)
		if err != nil {
			return apis.NewBadRequestError("failed to search members", err)
		}

		file, err := os.CreateTemp("", "setip_export_*.xlsx")
		if err != nil {
			return apis.NewBadRequestError("failed to create temporary file", err)
		}
		defer os.Remove(file.Name())

		err = service.Export(file, members, searchParams.Columns)
		if err != nil {
			return apis.NewBadRequestError("failed to export excel of members", err)
		}

		return ctx.Attachment(file.Name(), "export.xlsx")
	}
}

func Activate(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		var data service.ActivateMemberRequest

		err := rest.CopyJsonBody(ctx.Request(), &data)
		if err != nil {
			return apis.NewBadRequestError("failed to copy data", err)
		}

		record, err := app.Dao().FindRecordById("members", id)
		if err != nil {
			return apis.NewBadRequestError("failed to find member", err)
		}

		err = service.ActivateMember(ctx, app, record, &data)
		if err != nil {
			return apis.NewBadRequestError("failed to activate member", err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}
}

func Deactivate(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.PathParam("id")

		ctx.Set("dao", app.Dao())

		record, err := app.Dao().FindRecordById("members", id)
		if err != nil {
			return apis.NewBadRequestError("failed to find member", err)
		}

		err = service.DeactivateMember(ctx, app, record)
		if err != nil {
			return apis.NewBadRequestError("failed to deactivate member", err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}
}

func Import(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Set("dao", app.Dao())

		// Retrieve the file from the form data
		file, err := ctx.FormFile("file")
		if err != nil {
			return apis.NewBadRequestError("Unable to retrieve file", err)
		}

		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "Unable to open file", err)
		}
		defer src.Close()

		membersByName, companiesByName, err := service.Import(ctx, app, src)
		if err != nil {
			return apis.NewBadRequestError(
				fmt.Sprintf("failed to import excel of members: %s", err),
				err,
			)
		}

		return ctx.JSON(http.StatusOK, map[string]any{
			"members":   membersByName,
			"companies": companiesByName,
		})
	}
}

func prepareMembersearch(dao *daos.Dao, searchParams *query.SearchParams) (*query.SearchParams, error) {
	// Chapter, if present, will override the address related filters.
	if searchParams.ChapterID != "" {
		chapter, found, err := chapterQuery.FindByID(dao, searchParams.ChapterID)
		if err != nil {
			return nil, fmt.Errorf("failed to find chapter: %w", err)
		}

		if !found {
			return nil, errors.New("chapter not found")
		}

		searchParams.AddressCityIDs = chapter.CityIDs
		searchParams.LegacyArea = chapter.RawCityQuery
	}

	return searchParams, nil
}
