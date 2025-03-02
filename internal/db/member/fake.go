package member

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/fragoulis/setip_v2/internal/db/address"
)

func NewRandom(dao *daos.Dao) (*Member, error) {
	return New(
		WithNextNo(dao),
		WithRandomPersonalDetails(),
		WithRandomAddress(dao),
	)
}

func WithNextNo(dao *daos.Dao) Option {
	return func(mem *Member) error {
		no, err := LastNo(dao)
		if err != nil {
			return err
		}

		mem.MemberNo = no + 1

		return nil
	}
}

func WithRandomPersonalDetails() Option {
	return func(mem *Member) error {
		// Generate a random person (common name and email)
		person := gofakeit.Person()

		// Generate a birthdate in past 50 years.
		birthdate, err := types.ParseDateTime(gofakeit.DateRange(
			time.Now().AddDate(-50, 0, 0),
			time.Now().AddDate(-18, 0, 0),
		))
		if err != nil {
			return err
		}

		mem.FirstName = person.FirstName
		mem.LastName = person.LastName
		mem.FatherName = gofakeit.FirstName()
		mem.Email = person.Contact.Email
		mem.Mobile = fmt.Sprintf("69%08d", gofakeit.Number(0, 99999999))
		mem.Birthdate = birthdate

		return nil
	}
}

func WithRandomAddress(dao *daos.Dao) Option {
	return func(mem *Member) error {
		street, err := address.RandomStreet(dao)
		if err != nil {
			return fmt.Errorf("random street: %w", err)
		}

		mem.AddressCityID = street.CityID
		mem.AddressStreetID = street.GetId()

		//nolint:perfsprint
		mem.AddressStreetNo = fmt.Sprintf("%d", gofakeit.Number(0, 100))

		return nil
	}
}
