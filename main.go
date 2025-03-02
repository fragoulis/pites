package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/fragoulis/setip_v2/cmd"
	addressAPI "github.com/fragoulis/setip_v2/internal/app/address/api"
	assembly "github.com/fragoulis/setip_v2/internal/app/assembly"
	chaptersAPI "github.com/fragoulis/setip_v2/internal/app/chapter/api"
	companyAPI "github.com/fragoulis/setip_v2/internal/app/company/api"
	memberAPI "github.com/fragoulis/setip_v2/internal/app/member/api"
	paymentAPI "github.com/fragoulis/setip_v2/internal/app/payment/api"
	dbCompany "github.com/fragoulis/setip_v2/internal/db/company"
	dbEmployment "github.com/fragoulis/setip_v2/internal/db/employment"
	dbMember "github.com/fragoulis/setip_v2/internal/db/member"
	dbPayment "github.com/fragoulis/setip_v2/internal/db/payment"
	dbSubscription "github.com/fragoulis/setip_v2/internal/db/subscription"
	_ "github.com/fragoulis/setip_v2/migrations"
	"github.com/fragoulis/setip_v2/ui"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(srvEvnt *core.ServeEvent) error {
		memberAPI.RegisterRoutes(srvEvnt, app)
		companyAPI.RegisterRoutes(srvEvnt, app)
		addressAPI.RegisterRoutes(srvEvnt, app)
		paymentAPI.RegisterRoutes(srvEvnt, app)
		chaptersAPI.RegisterRoutes(srvEvnt, app)

		// serves static files from the provided public dir (if exists)
		srvEvnt.Router.GET("/*", apis.StaticDirectoryHandler(ui.DistDirFS, true))

		return nil
	})

	dbMember.RegisterCallbacks(app)
	dbCompany.RegisterCallbacks(app)
	assembly.RegisterCallbacks(app)
	dbEmployment.RegisterCallbacks(app)
	dbPayment.RegisterCallbacks(app)
	dbSubscription.RegisterCallbacks(app)

	app.RootCmd.AddCommand(cmd.NewUserCommand(app))
	app.RootCmd.AddCommand(cmd.NewSeedCommand(app))
	app.RootCmd.AddCommand(cmd.NewClearCommand(app))
	app.RootCmd.AddCommand(cmd.NewImportCommand(app))
	app.RootCmd.AddCommand(cmd.NewIssueCommand(app))
	app.RootCmd.AddCommand(cmd.NewMiscCommand(app))
	app.RootCmd.AddCommand(cmd.NewGenerateCommand(app))
	app.RootCmd.AddCommand(cmd.NewExportCommand(app))

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
