package company

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/daos"

	"github.com/fragoulis/setip_v2/internal/db/address"
)

func NewRandom(dao *daos.Dao, withBranch bool) (*Company, error) {
	return New(
		WithRandomName(),
		WithRandomBranch(withBranch),
		WithRandomAddress(dao),
	)
}

func WithRandomName() Option {
	return func(c *Company) error {
		c.Name = gofakeit.Company()

		return nil
	}
}

func WithRandomBranch(withBranch bool) Option {
	return func(c *Company) error {
		if withBranch {
			c.Branch = gofakeit.Blurb()
		}

		return nil
	}
}

func WithRandomAddress(dao *daos.Dao) Option {
	return func(comp *Company) error {
		street, err := address.RandomStreet(dao)
		if err != nil {
			return fmt.Errorf("random street: %w", err)
		}

		comp.AddressCityID = street.CityID
		comp.AddressStreetID = street.GetId()

		//nolint:perfsprint
		comp.AddressStreetNo = fmt.Sprintf("%d", gofakeit.Number(0, 100))

		return nil
	}
}
