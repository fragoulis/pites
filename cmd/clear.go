package cmd

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewClearCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Removes the data. Use with caution.",
	}

	command.AddCommand(clearAllCommand(app))
	command.AddCommand(clearAssembliesCommand(app))
	command.AddCommand(clearAuditlogsCommand(app))
	command.AddCommand(clearCompaniesCommand(app))
	command.AddCommand(clearEmploymentsCommand(app))
	command.AddCommand(clearMembersCommand(app))
	command.AddCommand(clearPaymentsCommand(app))
	command.AddCommand(clearSubscriptionsCommand(app))

	return command
}

func clearAllCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "all",
		Example:      "clear all",
		Short:        "Removes all data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			tables := []string{
				"assemblies",
				"auditlogs",
				"companies",
				"employments",
				"members",
				"payments",
				"subscriptions",
			}

			for _, table := range tables {
				if _, err := app.Dao().DB().
					NewQuery("delete from " + table).
					Execute(); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func clearAssembliesCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "assemblies",
		Example:      "clear assemblies",
		Short:        "Removes all assembly data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from assemblies").
				Execute()

			return err
		},
	}
}

func clearAuditlogsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "auditlogs",
		Example:      "clear auditlogs",
		Short:        "Removes all auditlog data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from auditlogs").
				Execute()

			return err
		},
	}
}

func clearCompaniesCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "companies",
		Example:      "clear companies",
		Short:        "Removes all company data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from companies").
				Execute()

			return err
		},
	}
}

func clearEmploymentsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "employments",
		Example:      "clear employments",
		Short:        "Removes all employment data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from employments").
				Execute()

			return err
		},
	}
}

func clearMembersCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "members",
		Example:      "clear members",
		Short:        "Removes all member data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from members").
				Execute()

			return err
		},
	}
}

func clearPaymentsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "payments",
		Example:      "clear payments",
		Short:        "Removes all payment data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from payments").
				Execute()

			return err
		},
	}
}

func clearSubscriptionsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "subscriptions",
		Example:      "clear subscriptions",
		Short:        "Removes all subscription data. Use with caution.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := app.Dao().DB().
				NewQuery("delete from subscriptions").
				Execute()

			return err
		},
	}
}
