package payment

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/daos"

	"github.com/fragoulis/setip_v2/internal/db/member"
)

const userID = "tatnjibskzq4sav"

func NewRandom(dao *daos.Dao) (*Payment, error) {
	return New(
		WithRandomMember(dao),
		WithRandomAmount(),
		WithRandomReceipt(),
		WithRandomUser(),
	)
}

func WithRandomMember(dao *daos.Dao) Option {
	return func(paym *Payment) error {
		member, err := member.Random(dao)
		if err != nil {
			return fmt.Errorf("random member: %w", err)
		}

		paym.MemberID = member.GetId()

		return nil
	}
}

func WithRandomAmount() Option {
	return func(s *Payment) error {
		s.Amount = 2 * gofakeit.Number(1, 5)

		return nil
	}
}

func WithRandomUser() Option {
	return func(s *Payment) error {
		s.CreatedByUserID = userID

		return nil
	}
}

//nolint:perfsprint
func WithRandomReceipt() Option {
	return func(s *Payment) error {
		s.ReceiptBlockNo = fmt.Sprintf("%d", gofakeit.Number(1, 500))
		s.ReceiptNo = fmt.Sprintf("%d", gofakeit.Number(1, 50))

		return nil
	}
}
