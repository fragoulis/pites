package payment

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*Payment)(nil)

type Payment struct {
	models.BaseModel

	MemberID        string         `db:"member_id"`
	Amount          int            `db:"amount_in_euros"`
	ReceiptBlockNo  string         `db:"receipt_block_no"`
	ReceiptNo       string         `db:"receipt_no"`
	IssueAt         types.DateTime `db:"issued_at"`
	LegacyGUID      string         `db:"legacy_guid"`
	CreatedByUserID string         `db:"created_by_user_id"`
}

func (m *Payment) TableName() string {
	return "payments"
}

type Option func(*Payment) error

func New(options ...Option) (*Payment, error) {
	model := &Payment{}

	for _, option := range options {
		err := option(model)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func Query(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Payment{})
}
