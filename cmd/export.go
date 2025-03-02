package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewExportCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "export",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("collection required")
			}

			collectionName := args[0]

			records, err := app.Dao().FindRecordsByFilter(
				collectionName, // collection
				"1=1",          // filter
				"",             // sort
				10000,          // limit
				0,              // offset
			)
			if err != nil {
				return fmt.Errorf("failed to fetch records: %w", err)
			}

			type Row map[string]any
			type Rows []Row

			rows := make(Rows, 0, len(records))

			for _, record := range records {
				rows = append(rows, record.ColumnValueMap())
			}

			raw, err := json.MarshalIndent(rows, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}

			fmt.Println(string(raw))

			return nil
		},
	}
}
