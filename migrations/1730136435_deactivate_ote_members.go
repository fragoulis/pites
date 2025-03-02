package migrations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/utils"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		ctx := &echo.DefaultContext{}
		ctx.Set("dao", dao)

		membersToDeactivate := map[string]string{
			// reducted
		}

		return dao.RunInTransaction(func(tx *daos.Dao) error {
			for code, name := range membersToDeactivate {
				memberNo, err := strconv.Atoi(code)
				if err != nil {
					return err
				}

				members, err := tx.FindRecordsByExpr("members", dbx.HashExp{"member_no": memberNo})
				if err != nil {
					return err
				}

				if len(members) == 0 {
					return fmt.Errorf("no member with code %s %s %d", code, name, memberNo)
				}

				member := members[0]

				name = utils.Normalize(name)
				lastName := member.GetString("last_name")
				firstName := member.GetString("first_name")
				fullName := lastName + " " + firstName

				if fullName != name {
					return fmt.Errorf("name and code do not match: %q %q = %q", code, name, fullName)
				}

				subscriptions, err := tx.FindRecordsByExpr(
					"subscriptions",
					dbx.HashExp{"member_id": member.GetId()},
				)
				if err != nil {
					return fmt.Errorf("failed to find subscription: %w", err)
				}

				if len(subscriptions) == 0 {
					fmt.Println("No subscriptions for", code, name)

					continue
				}

				activeSubscriptions := []*models.Record{}

				for _, subscription := range subscriptions {
					if subscription.GetDateTime("end_date").IsZero() {
						activeSubscriptions = append(activeSubscriptions, subscription)
					}
				}

				if len(activeSubscriptions) == 0 {
					fmt.Println("No active subscriptions for", code, name)

					continue
				}

				if len(activeSubscriptions) > 1 {
					fmt.Println("More than one active subscription for", code, name)
				}

				for _, subscription := range subscriptions {
					subscription.Set("active", false)
					subscription.Set("end_date", time.Now())

					if err := tx.SaveRecord(subscription); err != nil {
						return err
					}

					fmt.Println("Deactivated", code, name)
				}
			}

			return nil
		})
	}, func(_ dbx.Builder) error {
		return nil
	})
}
