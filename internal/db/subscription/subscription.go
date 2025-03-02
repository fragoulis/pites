package subscription

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*Subscription)(nil)

type Subscription struct {
	models.BaseModel

	MemberID  string         `db:"member_id"`
	Active    bool           `db:"active"`
	FeePaid   bool           `db:"fee_paid"`
	StartDate types.DateTime `db:"start_date"`
	EndDate   types.DateTime `db:"end_date"`
}

func (m *Subscription) TableName() string {
	return "subscriptions"
}

type Option func(*Subscription) error

func New(options ...Option) (*Subscription, error) {
	model := &Subscription{}

	for _, option := range options {
		err := option(model)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func Query(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Subscription{})
}
