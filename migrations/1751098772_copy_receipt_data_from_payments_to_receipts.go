package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		payments, err := dao.FindRecordsByFilter(
			"payments",
			"1=1",
			"-created",
			0,
			0,
		)
		if err != nil {
			return fmt.Errorf("failed to load payments: %w", err)
		}

		collection, err := dao.FindCollectionByNameOrId("receipts")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}

		err = dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, payment := range payments {
				receipt := models.NewRecord(collection)

				receipt.Set("member_id", payment.Get("member_id"))
				receipt.Set("amount_in_euros", payment.Get("amount_in_euros"))
				receipt.Set("block_no", payment.Get("receipt_block_no"))
				receipt.Set("receipt_no", payment.Get("receipt_no"))
				receipt.Set("issued_at", payment.Get("issued_at"))
				receipt.Set("comments", payment.Get("comments"))
				receipt.Set("created_by_user_id", payment.Get("created_by_user_id"))

				err = tx.Save(receipt)
				if err != nil {
					return fmt.Errorf("failed to create receipt record: %w", err)
				}

				payment.Set("receipt_id", receipt.GetId())

				err = tx.Save(payment)
				if err != nil {
					return fmt.Errorf("failed to update payment record: %w", err)
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to move receipts: %w", err)
		}

		return nil
	}, func(db dbx.Builder) error {
		_, err := db.
			NewQuery("UPDATE payments SET receipt_id = NULL").
			Execute()
		if err != nil {
			return fmt.Errorf("failed to empty receipt id in payments: %w", err)
		}

		_, err = db.
			NewQuery("DELETE FROM receipts").
			Execute()
		if err != nil {
			return fmt.Errorf("failed to truncate receipts: %w", err)
		}

		return nil
	})
}
