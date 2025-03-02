package cmd

import (
	"errors"

	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"
)

const defaultPassword = "password"

// NewUserCommand creates and returns new command for managing
// user accounts (create).
func NewUserCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "user",
		Short: "Manages user accounts",
	}

	command.AddCommand(userCreateCommand(app))

	return command
}

func userCreateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		Example:      "user create name",
		Short:        "Creates a new user account",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("missing username")
			}

			name := args[0]
			email := name + "@setip.gr"
			password := defaultPassword
			role := "user"

			if name == "" {
				return errors.New("missing username")
			}

			if email == "" || is.EmailFormat.Validate(email) != nil {
				return errors.New("missing or invalid email address")
			}

			if len(password) < 8 {
				return errors.New("the password must be at least 8 chars long")
			}

			collection, err := app.Dao().FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}

			record := models.NewRecord(collection)

			form := forms.NewRecordUpsert(app, record)

			err = form.LoadData(map[string]any{
				"username":        name,
				"name":            name,
				"email":           email,
				"role":            role,
				"password":        password,
				"passwordConfirm": password,
			})
			if err != nil {
				return err
			}

			if err := form.Submit(); err != nil {
				return err
			}

			color.Green("Successfully created new user: %s!", name)

			return nil
		},
	}

	return command
}
