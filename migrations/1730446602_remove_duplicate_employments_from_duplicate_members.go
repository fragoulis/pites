package migrations

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

// This migration fixes an issue that is the result
// of members being duplicate in the old database and
// our bad handling of the import.
// We created two employment records for each member
// record.

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		members, err := dao.FindRecordsByFilter(
			"members",
			"1=1",
			"",
			0,
			0,
		)
		if err != nil {
			return fmt.Errorf("failed to load members: %w", err)
		}

		return dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, member := range members {
				employments, err := dao.FindRecordsByFilter(
					"employments",
					"end_date = '' && member_id = {:member_id}",
					"-created",
					0,
					0,
					dbx.Params{"member_id": member.GetId()},
				)
				if err != nil {
					return fmt.Errorf("failed to load employments for member %q: %w", member.GetId(), err)
				}

				if len(employments) < 2 {
					continue
				}

				// Some have simply the same company twice
				companyIDs := map[string]bool{}
				for _, employment := range employments {
					companyIDs[employment.GetString("company_id")] = true
				}

				if len(companyIDs) == 1 {
					if err := tx.DeleteRecord(employments[0]); err != nil {
						return fmt.Errorf("failed to delete duplicate employment record: %w", err)
					}

					continue
				}

				// If all have a legacy guid, they are the result of
				// duplicate members and bad import on my side.
				// We end the oldest one.
				allWithLegacyGUID := true

				for _, employment := range employments {
					if employment.GetString("legacy_guid") == "" {
						allWithLegacyGUID = false

						break
					}
				}

				if !allWithLegacyGUID {
					continue
				}

				employmentToEnd := employments[1:]

				for _, employment := range employmentToEnd {
					employment.Set("end_date", time.Now())

					if err := tx.SaveRecord(employment); err != nil {
						return fmt.Errorf(
							"failed to end employment %q for member %q: %w",
							employment.GetId(),
							member.GetId(),
							err,
						)
					}
				}
			}

			// Delete this particular employment rj50xfu2aa4d4s0 which falls
			// under the category duplicate with the same company, but
			// fails to be picked by the query because there is also another employment
			// present.

			ids := []string{
				"rj50xfu2aa4d4s0",
				"y8hfgaep66p91hw",
				"e4uzs4lvnxfet7q",
				"j3adhibbu2ndajm",
				"u6waa9helhfb57w",
				"g69fqvgrt9miy0r",
			}

			for _, id := range ids {
				_, err := db.
					NewQuery(fmt.Sprintf("delete from employments where id='%s'", id)).
					Execute()
				if err != nil {
					return err
				}
			}

			return nil
		})
	}, func(_ dbx.Builder) error {
		return nil
	})
}
