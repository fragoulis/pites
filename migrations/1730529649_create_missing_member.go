package migrations

import (
	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(_ dbx.Builder) error {
		// reducted
		return nil
	}, func(_ dbx.Builder) error {
		return nil
	})
}
