package cmd

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewMiscCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use: "misc",
	}

	command.AddCommand(resaveMembers(app))

	return command
}

func resaveMembers(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "resave",
		Example:      "misc resave",
		Short:        "",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			total := 5300
			limit := 10
			offset := 0
			updatedTotal := 0

			for {
				records, err := app.Dao().FindRecordsByFilter(
					"members", // collection
					"1=1",     // filter
					"",        // sort
					limit,     // limit
					offset,    // offset
				)
				if err != nil {
					return fmt.Errorf("failed to fetch records: %w", err)
				}

				updatedNow := 0

				for _, record := range records {
					err := app.Dao().Save(record)
					if err != nil {
						return fmt.Errorf("failed to save: %w", err)
					}

					updatedTotal++
					updatedNow++
				}

				offset += limit

				if updatedNow == 0 || updatedTotal > total {
					break
				}

				fmt.Printf("%d/%d\n", updatedTotal, total)
			}

			return nil
		},
	}
}
