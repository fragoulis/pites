package member

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*Member)(nil)

type Member struct {
	models.BaseModel

	MemberNo          int            `db:"member_no"`
	FullName          string         `db:"full_name"`
	FirstName         string         `db:"first_name"`
	LastName          string         `db:"last_name"`
	FatherName        string         `db:"father_name"`
	AddressCityID     string         `db:"address_city_id"`
	AddressStreetID   string         `db:"address_street_id"`
	AddressStreetNo   string         `db:"address_street_no"`
	Email             string         `db:"email"`
	Mobile            string         `db:"mobile"`
	Phone             string         `db:"phone"`
	Birthdate         types.DateTime `db:"birthdate"`
	IDCardNumber      string         `db:"id_card_number"`
	SocialSecurityNum string         `db:"social_security_num"`
	OtherUnion        bool           `db:"other_union"`
	LegacyGUID        string         `db:"legacy_guid"`
	LegacyAddress     string         `db:"legacy_address"`
	SearchTerms       string         `db:"search_terms"`
}

func (m *Member) TableName() string {
	return "members"
}

type Option func(*Member) error

func New(options ...Option) (*Member, error) {
	model := &Member{}

	for _, option := range options {
		err := option(model)
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func WithNo(memberNo int) Option {
	return func(m *Member) error {
		m.MemberNo = memberNo

		return nil
	}
}

func Query(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Member{})
}

func Random(dao *daos.Dao) (*Member, error) {
	model := &Member{}

	err := Query(dao).
		OrderBy("RANDOM()").
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}

func FindByID(dao *daos.Dao, id string) (*Member, error) {
	model := &Member{}

	err := Query(dao).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model, nil
}

func LastNo(dao *daos.Dao) (int, error) {
	models := []*Member{}

	err := Query(dao).
		OrderBy("member_no desc").
		Limit(1).
		All(&models)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(models) == 0 {
		return 1, nil
	}

	return models[0].MemberNo, nil
}
