package company

import (
	"errors"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/db/auditlog"
	"github.com/fragoulis/setip_v2/internal/utils"
)

func setSearchTermsForCompany(model *Company) {
	model.Name = utils.Normalize(model.Name)
}

func setSearchTermsForRecord(model *models.Record) {
	model.Set("name", utils.Normalize(model.GetString("name")))
}

func setSearchTerms(evt *core.ModelEvent) error {
	company, ok := evt.Model.(*Company)
	if ok {
		setSearchTermsForCompany(company)

		return nil
	}

	model, ok := evt.Model.(*models.Record)
	if ok {
		setSearchTermsForRecord(model)

		return nil
	}

	return errors.New("failed to cast model to company")
}

func RegisterCallbacks(app *pocketbase.PocketBase) {
	app.OnModelBeforeCreate("companies").Add(setSearchTerms)
	app.OnModelBeforeUpdate("companies").Add(setSearchTerms)
	auditlog.RegisterCallbacks(app, "companies")
}
