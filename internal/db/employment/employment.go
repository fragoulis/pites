package employment

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*Employment)(nil)

type Employment struct {
	models.BaseModel

	MemberID  string         `db:"member_id"  json:"member_id"`
	CompanyID string         `db:"company_id" json:"company_id"`
	StartDate types.DateTime `db:"start_date" json:"start_date"`
	EndDate   types.DateTime `db:"end_date"   json:"end_date"`
}

func (m *Employment) TableName() string {
	return "employments"
}

type Option func(*Employment) error

func New(options ...Option) (*Employment, error) {
	model := &Employment{}

	for _, option := range options {
		err := option(model)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func Query(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Employment{})
}

func FindByMemberID(dao *daos.Dao, memberIDs []string) ([]*Employment, error) {
	records := []*Employment{}

	err := Query(dao).
		Where(dbx.In("member_id", list.ToInterfaceSlice(memberIDs)...)).
		OrderBy("start_date DESC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return records, nil
}
