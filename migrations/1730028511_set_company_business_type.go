package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		businessTypeIDByRawName := map[string]string{
			// reducted
		}

		businessTypeNameByCompanyNameRaw := map[string]string{
			// reducted
		}

		companies, err := dao.FindRecordsByFilter(
			"companies",    // collection
			"parent_id=''", // filter
			"",             // sort
			1000,           // limit
			0,              // offset
		)
		if err != nil {
			return err
		}

		updated := 0
		skipped := 0

		for _, company := range companies {
			companyName := company.GetString("name")

			typeNameRaw, ok := businessTypeNameByCompanyNameRaw[companyName]
			if !ok {
				skipped++

				continue
			}

			typeID, ok := businessTypeIDByRawName[typeNameRaw]
			if !ok {
				return fmt.Errorf("failed to find id for business type %q", typeNameRaw)
			}

			company.Set("business_type_id", typeID)

			if err := dao.SaveRecord(company); err != nil {
				return err
			}

			updated++
		}

		return nil
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery("update companies set business_type_id = ''").Execute()

		return err
	})
}
