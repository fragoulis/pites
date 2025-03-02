package member

import (
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/db/auditlog"
	"github.com/fragoulis/setip_v2/internal/utils"
)

func setAttributesForModel(model *Member) {
	model.FirstName = utils.Normalize(model.FirstName)
	model.LastName = utils.Normalize(model.LastName)
	model.FatherName = utils.Normalize(model.FatherName)
	model.Email = utils.Normalize(model.Email)
	model.FullName = fmt.Sprintf("%s %s", model.FirstName, model.LastName)
}

func setAttributesForRecord(model *models.Record) {
	model.Set("first_name", utils.Normalize(model.GetString("first_name")))
	model.Set("last_name", utils.Normalize(model.GetString("last_name")))
	model.Set("father_name", utils.Normalize(model.GetString("father_name")))
	model.Set("email", utils.Normalize(model.GetString("email")))
	model.Set("full_name", fmt.Sprintf("%s %s",
		model.GetString("first_name"),
		model.GetString("last_name"),
	))
}

func setAttributes(evt *core.ModelEvent) error {
	company, ok := evt.Model.(*Member)
	if ok {
		setAttributesForModel(company)

		return nil
	}

	model, ok := evt.Model.(*models.Record)
	if ok {
		setAttributesForRecord(model)

		return nil
	}

	return errors.New("failed to cast model to member")
}

func RegisterCallbacks(app *pocketbase.PocketBase) {
	app.OnModelBeforeCreate("members").Add(setAttributes)
	app.OnModelBeforeUpdate("members").Add(setAttributes)
	auditlog.RegisterCallbacks(app, "members")
}
