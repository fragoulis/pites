package migrations

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/list"
)

// In this migration, we go through some updates once more.
// In the initial update we didn't pay attention and we failed
// to end current employments.
// For each of there, if there is another open employment, we
// end it.

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		codes := []int{
			// reducted
		}

		members, err := dao.FindRecordsByExpr(
			"members",
			dbx.In("member_no", list.ToInterfaceSlice(codes)...),
		)
		if err != nil {
			return fmt.Errorf("failed to load members: %w", err)
		}

		return dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, member := range members {
				employments, err := tx.FindRecordsByFilter(
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

				// End all employments with legacy guid

				for _, employment := range employments {
					if employment.GetString("legacy_guid") == "" {
						continue
					}

					if !employment.GetDateTime("end_date").IsZero() {
						continue
					}

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

			return nil
		})
	}, func(_ dbx.Builder) error {
		return nil
	})
}
