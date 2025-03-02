package migrations

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		members, err := dao.FindRecordsByFilter(
			"members",
			"legacy_record.RegisterDate=\"\" && legacy_record.IsActive~true",
			"",
			0,
			0,
		)
		if err != nil {
			return fmt.Errorf("failed to load members: %w", err)
		}

		subscriptionsCollection, err := dao.FindCollectionByNameOrId("subscriptions")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}

		return dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, member := range members {
				if member.GetString("legacy_record") == "" {
					return fmt.Errorf("member missing legacy record: %s", member.GetId())
				}

				var legacyRecord struct {
					//nolint:tagliatelle
					ValidThrough string `json:"ValidThrough"`
				}

				err := json.Unmarshal([]byte(member.GetString("legacy_record")), &legacyRecord)
				if err != nil {
					return fmt.Errorf("failed to unmarshal: %w", err)
				}

				// Get the validThrough date from the legacy record.
				startDate, err := time.Parse(time.DateOnly, strings.Split(legacyRecord.ValidThrough, " ")[0])
				if err != nil {
					return fmt.Errorf("failed to parse date: %w", err)
				}

				if startDate.IsZero() {
					startDate = time.Now()
				}

				// Get member's subscriptions
				subscriptions, err := dao.FindRecordsByExpr(
					"subscriptions",
					dbx.HashExp{"member_id": member.GetId()},
				)
				if err != nil {
					return fmt.Errorf("failed to find member subscriptions: %w", err)
				}

				// If no subscriptions found create one
				if len(subscriptions) == 0 {
					subscription := models.NewRecord(subscriptionsCollection)
					subscription.Set("member_id", member.GetId())
					subscription.Set("fee_paid", true)
					subscription.Set("active", true)
					subscription.Set("start_date", startDate)
					subscription.Set("legacy_guid", member.GetString("legacy_guid"))

					if err := tx.SaveRecord(subscription); err != nil {
						return fmt.Errorf("failed to create subscription %q: %w", member.GetId(), err)
					}

					continue
				}

				foundActive := false

				for _, subscription := range subscriptions {
					// If the member has an active subscription just make
					// sure it is also active and continue
					if subscription.GetDateTime("end_date").IsZero() {
						subscription.Set("start_date", startDate)
						subscription.Set("fee_paid", true)
						subscription.Set("active", true)

						// fmt.Println("Updating active subscription for", member.GetString("full_name"))

						if err := tx.SaveRecord(subscription); err != nil {
							return fmt.Errorf("failed to update subscription active %q: %w", member.GetId(), err)
						}

						foundActive = true

						break
					}
				}

				if foundActive {
					continue
				}

				// If no active subscription was found take the first one
				// reactivate
				if len(subscriptions) > 1 {
					return fmt.Errorf("member has more than one subscription %s", member.GetId())
				}

				subscription := subscriptions[0]
				subscription.Set("start_date", startDate)
				subscription.Set("fee_paid", true)
				subscription.Set("active", true)

				if err := tx.SaveRecord(subscription); err != nil {
					return fmt.Errorf("failed to update subscription %q: %w", member.GetId(), err)
				}
			}

			return nil
		})
	}, func(_ dbx.Builder) error {
		return nil
	})
}
