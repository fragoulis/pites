package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		members, err := dao.FindRecordsByFilter(
			"members",
			"company_id = ''",
			"",
			0,
			0,
		)
		if err != nil {
			return fmt.Errorf("failed to load members: %w", err)
		}

		var errMemberNo []string

		err = dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, member := range members {
				// For each member, load the open employments.
				employments, err := dao.FindRecordsByFilter(
					"employments",
					"end_date = '' && member_id = {:member_id}",
					"",
					0,
					0,
					dbx.Params{"member_id": member.GetId()},
				)
				if err != nil {
					return fmt.Errorf("failed to load employments for member %q: %w", member.GetId(), err)
				}

				if len(employments) == 0 {
					continue
				}

				if len(employments) > 1 {
					errMemberNo = append(errMemberNo, member.GetString("member_no"))

					continue
				}

				employment := employments[0]

				// Pick the branch id if present, otherwise the company id.
				// Both IDs are pointing to a company record row.
				companyID := employment.GetString("branch_id")
				if companyID == "" {
					companyID = employment.GetString("company_id")
				}

				if companyID == "" {
					return fmt.Errorf("company id cannot be empty: %q %q", member.GetId(), employment.GetId())
				}

				member.Set("company_id", companyID)

				if err := tx.SaveRecord(member); err != nil {
					return fmt.Errorf("failed to set company to member %q: %w", member.GetId(), err)
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("transcation failed: %w", err)
		}

		if len(errMemberNo) > 0 {
			return fmt.Errorf("found members with multiple employments: %v", errMemberNo)
		}

		return nil
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery("update members set company_id = ''").Execute()
		if err != nil {
			return err
		}

		return nil
	})
}
