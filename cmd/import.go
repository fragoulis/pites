//nolint
package cmd

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	memberQuery "github.com/fragoulis/setip_v2/internal/app/member/query"
	"github.com/fragoulis/setip_v2/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

type MemberJSONRow struct {
	Code              string `json:"Code"`
	FirstName         string `json:"FirstName"`
	LastName          string `json:"LastName"`
	FullName          string `json:"FullName"`
	FatherName        string `json:"FatherName"`
	RegisterDate      string `json:"RegisterDate"`
	ValidThrough      string `json:"ValidThrough"`
	BirthYear         int    `json:"BirthYear"`
	BirthDate         string `json:"BirthDate"`
	Unemployed        bool   `json:"Unemployed"`
	IsActive          bool   `json:"IsActive"`
	RecordGUID        string `json:"RecordGuid"`
	Phone             string `json:"Phone"`
	Mobile            string `json:"Mobile"`
	Email             string `json:"EMail"`
	CompanyGUID       string `json:"CompanyGuid"`
	CompanyName       string `json:"CompanyName"`
	CompanyBranch     string `json:"CompanyBranch"`
	Address           string `json:"Address"`
	PostCode          string `json:"PostCode"`
	Area              string `json:"Area"`
	City              string `json:"City"`
	IDCardNumber      string `json:"IDCardNumber"`
	IDCardAuthority   string `json:"IDCardAuthority"`
	OtherUnion        bool   `json:"OtherUnion"`
	Specialty         string `json:"Specialty"`
	Education         string `json:"Education"`
	SocialSecurityNum string `json:"SocialSecurityNum"`
}

type CompanyJSONRow struct {
	FullName    string `json:"FullName"`
	MainAddress string `json:"MainAddress"`
	Branch      string `json:"Branch"`
	Phone1      string `json:"Phone1"`
	Email       string `json:"EMail"`
	WebSite     string `json:"WebSite"`
	IsActive    bool   `json:"IsActive"`
	RecordGUID  string `json:"RecordGuid"`
}

type ImportReport struct {
	RecordsCreated int `json:"created"`
	RecordsUpdated int `json:"updated"`
	Duplicates     int `json:"duplicates"`
}

func (r *ImportReport) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		fmt.Println("failed to marshal report")

		return fmt.Sprintf(`
RecordsCreated: %d
RecordsUpdated: %d
Duplicates: %d
`,
			r.RecordsCreated,
			r.RecordsUpdated,
			r.Duplicates,
		)
	}

	return string(bytes)
}

//nolint:gochecknoglobals
var importReport = ImportReport{}

func NewImportCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "import",
		Short: "Imports data from legacy format.",
	}

	command.AddCommand(importCompaniesCommand(app))
	command.AddCommand(importMembersCommand(app))
	command.AddCommand(importSubscriptionsCommand(app))
	command.AddCommand(importEmploymentsCommand(app))
	command.AddCommand(importPaymentsCommand(app))
	command.AddCommand(importDeletionsCommand(app))

	return command
}

//nolint:varnamelen
func readJSONFile(file string, v any) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		//nolint:wrapcheck
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	//nolint:wrapcheck
	return json.Unmarshal(byteValue, v)
}

func normalizePhoneNumber(phoneNumber string) string {
	return strings.TrimSpace(phoneNumber)
}

func findOrInitializeByLegacyID(dao *daos.Dao, collection *models.Collection, guid string) (*models.Record, error) {
	var record *models.Record
	var err error

	if guid != "" {
		record, err = dao.FindFirstRecordByData(collection.GetId(), "legacy_guid", guid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to execute query for %s %q: %w", collection.GetId(), guid, err)
		}
	}

	if record == nil {
		record = models.NewRecord(collection)
	}

	return record, nil
}

func findBy(dao *daos.Dao, collection *models.Collection, key, value string) (*models.Record, error) {
	record, err := dao.FindFirstRecordByData(collection.GetId(), key, value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			//nolint:wrapcheck
			return nil, err
		}
	}

	return record, nil
}

func findCompanyBranch(
	dao *daos.Dao,
	collection *models.Collection,
	companyID string,
	name string,
) (*models.Record, error) {
	record := &models.Record{}

	err := dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{inflector.Columnify("parent_id"): companyID}).
		AndWhere(dbx.HashExp{inflector.Columnify("name"): name}).
		Limit(1).
		One(record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			//nolint:wrapcheck
			return nil, err
		}
	}

	return record, nil
}

func findMemberByGUID(dao *daos.Dao, value string) (*models.Record, error) {
	rec, err := dao.FindFirstRecordByData("members", "legacy_guid", value)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to execute query for %s %q: %w", "members", value, err)
	}

	if rec == nil {
		rec, err = dao.FindFirstRecordByData("members", "legacy_guid_2", value)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to execute query for %s %q: %w", "members", value, err)
		}
	}

	return rec, nil
}

func findMemberByCode(dao *daos.Dao, value string) (*models.Record, error) {
	rec, err := dao.FindFirstRecordByData("members", "member_no", value)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to execute query for %s %q: %w", "members", value, err)
	}

	return rec, nil
}

func findCompanyByGUIDOrName(dao *daos.Dao, value string) (*models.Record, error) {
	rec, err := dao.FindFirstRecordByData("companies", "legacy_guid", value)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to execute query for %s %q: %w", "companies", value, err)
	}

	if rec == nil {
		rec, err = dao.FindFirstRecordByData("companies", "name", value)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to execute query for %s %q: %w", "companies", value, err)
		}
	}

	return rec, nil
}

//nolint:funlen,gocognit
func importCompaniesCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "companies",
		Example:      "import companies",
		Short:        "Imports companies from json",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]

			var rows []CompanyJSONRow

			err := readJSONFile(filename, &rows)
			if err != nil {
				return err
			}

			collection, err := app.Dao().FindCollectionByNameOrId("companies")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			unknownPhone := utils.Normalize("Άγνωστο")
			unknownAddress := utils.Normalize("Άγνωστη")

			return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
				for _, row := range rows {
					var record *models.Record

					row.FullName = utils.Normalize(row.FullName)

					if row.FullName == "EVALUE" {
						row.FullName = "E-VALUE"
					}

					row.Branch = utils.Normalize(row.Branch)
					row.MainAddress = utils.Normalize(row.MainAddress)
					row.Phone1 = normalizePhoneNumber(row.Phone1)

					if row.Phone1 != "0000000000" || row.Phone1 == unknownPhone {
						row.Phone1 = ""
					}

					if row.MainAddress == unknownAddress {
						row.MainAddress = ""
					}

					if row.RecordGUID != "" {
						// Try to find duplicate by legacy guid
						record, err = findBy(tx, collection, "legacy_guid", row.RecordGUID)
						if err != nil {
							return fmt.Errorf("failed to find by company legacy_guid for %q: %w", row.RecordGUID, err)
						}
					}

					if record == nil {
						// Try to find duplicate by name
						record, err = findBy(tx, collection, "name", row.FullName)
						if err != nil {
							return fmt.Errorf("failed to find by company name for %q: %w", row.FullName, err)
						}
					}

					if record == nil {
						record = models.NewRecord(collection)
						record.Set("name", row.FullName)
						record.Set("phone", row.Phone1)
						record.Set("email", utils.Normalize(row.Email))
						record.Set("website", utils.Normalize(row.WebSite))
						record.Set("legacy_address", row.MainAddress)

						// log.Printf("Created: %q\n", row.FullName)

						importReport.RecordsCreated++
					} else {
						if row.MainAddress != "" {
							record.Set("legacy_address", row.MainAddress)
						}

						if row.Phone1 != "" {
							record.Set("phone", row.Phone1)
						}

						if record.GetString("phone") == "0000000000" {
							record.Set("phone", "")
						}

						// log.Printf("Updated: %q\n", row.FullName)

						importReport.RecordsUpdated++
					}

					if row.RecordGUID != "" {
						rawRecord, err := json.Marshal(row)
						if err != nil {
							return fmt.Errorf("failed to marshal row: %w", err)
						}

						record.Set("legacy_guid", row.RecordGUID)
						record.Set("legacy_record", rawRecord)
					}

					if err := tx.SaveRecord(record); err != nil {
						return fmt.Errorf("failed to save company %q: %w", row.FullName, err)
					}

					if row.Branch == "" {
						continue
					}

					branchRecord, err := findCompanyBranch(tx, collection, record.GetId(), row.Branch)
					if err != nil {
						return fmt.Errorf("failed to find company branch for %q: %w", record.GetId(), err)
					}

					if branchRecord != nil {
						continue
					}

					branchRecord = models.NewRecord(collection)
					branchRecord.Set("parent_id", record.GetId())
					branchRecord.Set("name", row.Branch)

					importReport.RecordsCreated++

					if err := tx.SaveRecord(branchRecord); err != nil {
						return fmt.Errorf("failed to save company branch %q-%q: %w", record.GetId(), row.Branch, err)
					}
				}

				return nil
			})
		},
	}
}

// Create members from the members.json dump.
//
//nolint:funlen,gocognit,cyclop
func importMembersCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "members",
		Example:      "import members",
		Short:        "Imports members from legacy format.",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]

			// We intentionally read in the entire row payload in order to save it as it is
			// to the database.

			var rows []MemberJSONRow

			err := readJSONFile(filename, &rows)
			if err != nil {
				return err
			}

			collection, err := app.Dao().FindCollectionByNameOrId("members")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
				for _, row := range rows {
					if row.Code == "000000" {
						log.Printf("skipped entry %q\n", row.RecordGUID)

						continue
					}

					// Keep the original raw record here.
					rawRecord, err := json.Marshal(row)
					if err != nil {
						return fmt.Errorf("failed to marshal row: %w", err)
					}

					row.FirstName = utils.Normalize(row.FirstName)
					row.LastName = utils.Normalize(row.LastName)
					row.FullName = utils.Normalize(row.FullName)
					row.FatherName = utils.Normalize(row.FatherName)
					row.Phone = normalizePhoneNumber(row.Phone)
					row.Mobile = normalizePhoneNumber(row.Mobile)
					row.Email = utils.Normalize(row.Email)
					row.Address = utils.Normalize(row.Address)
					row.PostCode = utils.Normalize(row.PostCode)
					row.Area = utils.Normalize(row.Area)
					row.City = utils.Normalize(row.City)
					row.IDCardNumber = utils.Normalize(row.IDCardNumber)
					row.IDCardAuthority = utils.Normalize(row.IDCardAuthority)
					row.Specialty = utils.Normalize(row.Specialty)
					row.Education = utils.Normalize(row.Education)
					row.SocialSecurityNum = utils.Normalize(row.SocialSecurityNum)

					if row.FatherName == "ΑΓΝΩΣΤΟ" {
						row.FatherName = ""
					}

					if row.IDCardNumber == "Α123456" || row.IDCardNumber == "ΑΓΝΩΣΤΟ" || row.IDCardNumber == "Α00" {
						row.IDCardNumber = ""
					}

					//
					// Find the correct record.
					// Instead of goind directly to find by guid, first try email and mobile
					// which are unique to people, and as a result will allow us to catch
					// duplicate member records.
					//

					// Try to find duplicate by full name
					// record, err := findBy(tx, collection, "full_name", row.FullName)
					// if err != nil {
					// 	return fmt.Errorf("failed to find by full_name for %q: %w", row.FullName, err)
					// }

					record, err := findBy(tx, collection, "member_no", row.Code)
					if err != nil {
						return fmt.Errorf("failed to find by member_no for %q: %w", row.FullName, err)
					}

					// Try to find duplicate by email
					// if record == nil && row.Email != "" {
					// 	record, err = findBy(tx, collection, "email", row.Email)
					// 	if err != nil {
					// 		return fmt.Errorf("failed to find by email for %q: %w", row.FullName, err)
					// 	}
					// }

					// Try to find duplicate by id card number
					// if record == nil && row.IDCardNumber != "" {
					// 	record, err = findBy(tx, collection, "id_card_number", row.IDCardNumber)
					// 	if err != nil {
					// 		return fmt.Errorf("failed to find by id_card_number for %q: %w", row.FullName, err)
					// 	}
					// }

					// Try to find duplicate by mobile
					// if record == nil && row.Mobile != "" {
					// 	record, err = findBy(tx, collection, "mobile", row.Mobile)
					// 	if err != nil {
					// 		return fmt.Errorf("failed to find by modile for %q: %w", row.FullName, err)
					// 	}

					// 	// CODE USED FOR INITIAL IMPORT
					// 	// There is at least one case where 2 different people have the same mobile
					// 	// and as a result we cannot confidently merge the 2 records.
					// 	// Before doing so, compare their names.

					// 	// if record != nil &&
					// 	// 	row.LastName[:4] != utils.Normalize(record.GetString("last_name"))[:4] {
					// 	// 	log.Printf(
					// 	// 		"members share mobile number %q %q %s\n",
					// 	// 		row.RecordGUID,
					// 	// 		record.GetString("legacy_guid"),
					// 	// 		row.Mobile,
					// 	// 	)

					// 	// 	record = nil
					// 	// }
					// }

					// CODE USED FOR INITIAL IMPORT
					// // If there is indeed a duplicate record:
					// if record != nil {
					// 	importReport.Duplicates++

					// 	if row.IsActive {
					// 		// Overwrite with latest version only if latest is active and the
					// 		// new record has a guid (is legacy).
					// 		if row.RecordGUID != "" {
					// 			record.Set("legacy_guid_2", record.GetString("legacy_guid"))
					// 		}

					// 		record.Set("legacy_record_2", record.GetString("legacy_record"))
					// 	} else {
					// 		// Otherwise, simply store the latest version raw and move on.
					// 		if row.RecordGUID != "" {
					// 			record.Set("legacy_guid_2", row.RecordGUID)
					// 		}
					// 		record.Set("legacy_record_2", rawRecord)

					// 		if err := tx.SaveRecord(record); err != nil {
					// 			return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
					// 		}

					// 		continue
					// 	}
					// } else {
					// 	record, err = findOrInitializeByLegacyID(tx, collection, row.RecordGUID)
					// 	if err != nil {
					// 		return err
					// 	}
					// }

					originalRecord := map[string]any{}

					if record == nil {
						record = models.NewRecord(collection)
						importReport.RecordsCreated++

						record.Set("member_no", row.Code)
						record.Set("first_name", row.FirstName)
						record.Set("last_name", row.LastName)
						record.Set("full_name", row.FullName)
					} else {
						importReport.RecordsUpdated++

						originalRecord = record.OriginalCopy().PublicExport()
					}

					//
					// Update or create the record
					//

					if row.RecordGUID != "" {
						record.Set("legacy_guid", row.RecordGUID)
						record.Set("legacy_record", rawRecord)
					}

					// Update info if present, since in case of duplicates,
					// this record is probably most recent.
					if row.Mobile != "" {
						record.Set("mobile", row.Mobile)
					}
					if row.Phone != "" {
						record.Set("phone", row.Phone)
					}
					if row.Email != "" {
						record.Set("email", row.Email)
					}
					if row.FatherName != "" {
						record.Set("father_name", row.FatherName)
					}
					if row.Specialty != "" {
						record.Set("specialty", row.Specialty)
					}
					if row.Education != "" {
						record.Set("education", row.Education)
					}
					if row.IDCardNumber != "" {
						record.Set("id_card_number", row.IDCardNumber)
					}
					if row.IDCardAuthority != "" {
						record.Set("id_card_authority", row.IDCardAuthority)
					}
					if row.SocialSecurityNum != "" {
						record.Set("social_security_num", row.SocialSecurityNum)
					}
					if row.Address != "" {
						record.Set("legacy_address", row.Address)
					}
					if row.Area != "" {
						record.Set("legacy_area", row.Area)
					}
					if row.City != "" {
						record.Set("legacy_city", row.City)
					}
					if row.PostCode != "" {
						record.Set("legacy_post_code", row.PostCode)
					}

					record.Set("other_union", row.OtherUnion)

					// We have added a birthdate field because the new members excel included that
					// info and we thought it would be a waste to lose that info.
					if row.BirthDate != "" {
						birthdate, err := types.ParseDateTime(row.BirthDate)
						if err != nil {
							return fmt.Errorf("failed to parse birtdate for %q: %w", row.RecordGUID, err)
						}
						record.Set("birthdate", birthdate)
					} else if row.BirthYear > 0 {
						birthdate, err := types.ParseDateTime(fmt.Sprintf("%d-01-01", row.BirthYear))
						if err != nil {
							return fmt.Errorf("failed to parse birthyear for %q: %w", row.RecordGUID, err)
						}
						record.Set("birthdate", birthdate)
					}

					// Create diff

					dirtyRecord := record.PublicExport()

					action := "updating"
					if record.IsNew() {
						action = "creating"
					}

					for key, originalValue := range originalRecord {
						dirtyValue := dirtyRecord[key]

						switch key {
						case "created", "updated", "legacy_record", "legacy_record_2", "other_union":
							continue
						case "birthdate":
							if originalValue.(types.DateTime) == dirtyValue.(types.DateTime) {
								continue
							}
						case "member_no":
							if originalValue.(float64) == dirtyValue.(float64) {
								continue
							}
						default:
							if originalValue.(string) == dirtyValue.(string) {
								continue
							}
						}

						fmt.Printf("%s %q.%q from %q to %q\n", action, row.FullName, key, originalValue, dirtyValue)
					}

					// if err := tx.SaveRecord(record); err != nil {
					// 	return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
					// }
				}

				return nil
			})
		},
	}
}

// Create subscription from the members.json dump.
// By choice, we create one subscription record for each member and
// if this member is not active, we immediately end the subscription
// with today's date.
//
//nolint:funlen,gocognit
func importSubscriptionsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "subscriptions",
		Example:      "import subscriptions",
		Short:        "Imports subscriptions from legacy format.",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]

			var rows []MemberJSONRow

			err := readJSONFile(filename, &rows)
			if err != nil {
				return err
			}

			collection, err := app.Dao().FindCollectionByNameOrId("subscriptions")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
				for _, row := range rows {
					if row.Code == "000000" {
						log.Printf("skipped entry %q\n", row.RecordGUID)

						continue
					}

					registerDate, err := time.Parse("2006-01-02T00:00:00", row.RegisterDate)
					if err != nil {
						return fmt.Errorf("failed to parse register date: %w", err)
					}

					validThrough, err := time.Parse("2006-01-02T00:00:00", row.ValidThrough)
					if err != nil {
						return fmt.Errorf("failed to parse register date: %w", err)
					}

					if registerDate.IsZero() {
						if row.IsActive {
							// log.Printf("member has no registration date %q\n", row.RecordGUID)
						}

						continue
					}

					// Skip dates before the union start.
					if registerDate.Before(time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC)) {
						// log.Printf("member has invalid registration date %q %q\n", row.RecordGUID, row.RegisterDate)

						continue
					}

					var member *models.Record

					if row.RecordGUID != "" {
						member, err = findMemberByGUID(tx, row.RecordGUID)
						if err != nil {
							return err
						}
					}

					if member == nil {
						// New records do not have a guid.
						member, err = findMemberByCode(tx, row.Code)
						if err != nil {
							return err
						}
					}

					if member == nil {
						return fmt.Errorf("failed to find member %q %q", row.Code, row.RecordGUID)
					}

					legacyGUID1 := member.GetString("legacy_guid")

					var record *models.Record

					if legacyGUID1 != "" {
						// See if there is already a subscription for that member by the legacy guid.
						record, err = findOrInitializeByLegacyID(tx, collection, legacyGUID1)
						if err != nil {
							return err
						}
					}

					if record != nil && (record.GetBool("active") && row.IsActive) {
						// If there is second record
						//   and they are both active,
						//   do not create a second subscription, rather merge them,

						// Get starting date from existing subscription.
						var existing struct {
							RegisterDate types.DateTime `json:"RegisterDate"`
						}

						existingRaw := member.GetString("legacy_record_2")

						err = json.Unmarshal([]byte(existingRaw), &existing)
						if err != nil {
							return fmt.Errorf("failed to unmarshal: %w", err)
						}

						if registerDate.After(existing.RegisterDate.Time()) {
							log.Printf("WARN: Two active subscriptions found. skipping. %q\n", row.RecordGUID)

							continue
						}

						log.Printf(
							"WARN: Two active subscriptions found. updating. %q %q %q\n",
							row.Code,
							row.RegisterDate,
							existing.RegisterDate,
						)

						record.Set("start_date", row.RegisterDate)
						if err := tx.SaveRecord(record); err != nil {
							return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
						}

						continue
					} else {
						record = models.NewRecord(collection)
					}

					record.Set("legacy_guid", legacyGUID1)
					record.Set("member_id", member.GetId())
					record.Set("start_date", row.RegisterDate)

					if row.IsActive {
						record.Set("active", true)
					} else {
						if validThrough.IsZero() {
							return fmt.Errorf("valid through cannot be empty for inactive: %q", row.Code)
						}

						record.Set("end_date", row.ValidThrough)
					}

					if err := tx.SaveRecord(record); err != nil {
						return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
					}
				}

				return nil
			})
		},
	}
}

// Create employments from the members.json dump.
//
//nolint:funlen,gocognit,cyclop
func importEmploymentsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "employments",
		Example:      "import employments",
		Short:        "Imports employments from legacy format.",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]

			var rows []MemberJSONRow

			err := readJSONFile(filename, &rows)
			if err != nil {
				return err
			}

			collection, err := app.Dao().FindCollectionByNameOrId("employments")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			companies, err := app.Dao().FindCollectionByNameOrId("companies")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
				for _, row := range rows {
					if row.Code == "000000" {
						log.Printf("skipped entry %q\n", row.RecordGUID)

						continue
					}

					if row.CompanyGUID == "220b65c8-f1b1-4878-9fdb-d168298cdc78" {
						row.Unemployed = true
					}

					// If the member is unemployed currently, skip the entire employment record.
					if row.Unemployed {
						continue
					}

					var member *models.Record

					if row.RecordGUID != "" {
						member, err = findMemberByGUID(tx, row.RecordGUID)
						if err != nil {
							return err
						}
					}

					if member == nil {
						// New records do not have a guid.
						member, err = findMemberByCode(tx, row.Code)
						if err != nil {
							return err
						}
					}

					if member == nil {
						return fmt.Errorf("failed to find member %q %q", row.Code, row.RecordGUID)
					}

					row.CompanyName = utils.Normalize(row.CompanyName)

					if row.CompanyName == "EVALUE" {
						row.CompanyName = "E-VALUE"
					}

					var company *models.Record

					// Try to find company either by guid or name. New records do not have a guid.
					// We also have removed some duplicates by hand.
					if row.CompanyGUID != "" {
						company, err = findCompanyByGUIDOrName(tx, row.CompanyGUID)
						if err != nil {
							return err
						}
					}

					if company == nil && row.CompanyName != utils.Normalize("ΑΝΕΡΓΟΣ") {
						company, err = findCompanyByGUIDOrName(tx, row.CompanyName)
						if err != nil {
							return err
						}

						if company == nil {
							return fmt.Errorf("failed to find company %q %q", row.Code, row.CompanyName)
						}
					}

					// Find the branch
					var branch *models.Record

					if company != nil && row.CompanyBranch != "" {
						branch, err = findCompanyBranch(tx, companies, company.GetId(), row.CompanyBranch)
						if err != nil {
							return err
						}

						if branch == nil {
							return fmt.Errorf("failed to find company branch %q %q", row.Code, row.CompanyBranch)
						}
					}

					var record *models.Record

					// Try to find by legacy guid if present
					if row.RecordGUID != "" {
						record, err = findOrInitializeByLegacyID(tx, collection, row.RecordGUID)
						if err != nil {
							return err
						}
					} else if company == nil {
						// If no company is present at this point, the member is now unemployed.
						// Find the record without end date in order to close it.
						// If none is found, return.
						rec := &models.Record{}
						err := tx.RecordQuery(collection).
							AndWhere(dbx.HashExp{inflector.Columnify("member_id"): member.GetId()}).
							AndWhere(dbx.HashExp{inflector.Columnify("end_date"): ""}).
							Limit(1).
							One(rec)
						if err != nil {
							if !errors.Is(err, sql.ErrNoRows) {
								return fmt.Errorf("failed to execute query for employment: %w", err)
							}

							continue
						}

						record = rec
					} else {
						// If company is found, but no branch is present, try to find the right
						// record.
						rec := &models.Record{}
						err := tx.RecordQuery(collection).
							AndWhere(dbx.HashExp{inflector.Columnify("member_id"): member.GetId()}).
							AndWhere(dbx.HashExp{inflector.Columnify("company_id"): company.GetId()}).
							AndWhere(dbx.HashExp{inflector.Columnify("end_date"): ""}).
							Limit(1).
							One(rec)
						if err != nil {
							if !errors.Is(err, sql.ErrNoRows) {
								return fmt.Errorf("failed to execute query for employment: %w", err)
							}
						} else {
							record = rec
						}
					}

					originalRecord := map[string]any{}

					if record == nil {
						record = models.NewRecord(collection)
					} else {
						originalRecord = record.OriginalCopy().PublicExport()
					}

					if record.IsNew() {
						record.Set("member_id", member.GetId())
						record.Set("company_id", company.GetId())
						record.Set("start_date", types.NowDateTime())

						importReport.RecordsCreated++

						fmt.Printf("New employment for %q\n", row.FullName)

						// Close any current open employment before creating a new one.
						currentEmployment := &models.Record{}

						err := tx.RecordQuery(collection).
							AndWhere(dbx.HashExp{inflector.Columnify("member_id"): member.GetId()}).
							AndWhere(dbx.HashExp{inflector.Columnify("end_date"): ""}).
							Limit(1).
							One(currentEmployment)
						if err != nil {
							if !errors.Is(err, sql.ErrNoRows) {
								return fmt.Errorf("failed to execute query for employment: %w", err)
							}
						} else {
							currentEmployment.Set("end_date", types.NowDateTime())
							// if err := tx.SaveRecord(currentEmployment); err != nil {
							// 	return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
							// }

							fmt.Printf("Closing employment for %q\n", row.FullName)
						}
					} else {
						importReport.RecordsUpdated++
					}

					if branch != nil {
						record.Set("branch_id", branch.GetId())
					}

					if company == nil {
						record.Set("end_date", types.NowDateTime())
					}

					if !record.IsNew() {
						dirtyRecord := record.PublicExport()
						for _, key := range []string{"member_id", "company_id", "branch_id", "start_date", "end_date"} {
							originalValue, originalOk := originalRecord[key]
							dirtyValue := dirtyRecord[key]

							switch key {
							case "start_date", "end_date":
								if !originalOk {
									originalValue = types.DateTime{}
								}
								if originalValue.(types.DateTime) == dirtyValue.(types.DateTime) {
									continue
								}
							default:
								if !originalOk {
									originalValue = ""
								}
								if originalValue.(string) == dirtyValue.(string) {
									continue
								}
							}
							fmt.Printf("Updating %q.%q from %q to %q\n", row.FullName, key, originalValue, dirtyValue)
						}
					}

					if err := tx.SaveRecord(record); err != nil {
						return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
					}
				}

				return nil
			})
		},
	}
}

//nolint:funlen
func importPaymentsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "payments",
		Example:      "import payments",
		Short:        "Imports payments from legacy format.",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]

			//nolint:tagliatelle
			var rows []struct {
				MemberCode          string         `json:"MemberCode"`
				RecordGUID          string         `json:"RecordGuid"`
				MemberGUID          string         `json:"MemberGuid"`
				Amount              float64        `json:"Amount"`
				IssueDate           types.DateTime `json:"IssueDate"`
				FromMonth           types.DateTime `json:"FromMonth"`
				ToMonth             types.DateTime `json:"ToMonth"`
				Months              int            `json:"Months"`
				ContainsRegisterFee bool           `json:"ContainsRegisterFee"`
			}

			err := readJSONFile(filename, &rows)
			if err != nil {
				return err
			}

			collection, err := app.Dao().FindCollectionByNameOrId("payments")
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			return app.Dao().RunInTransaction(func(tx *daos.Dao) error {
				for _, row := range rows {
					if row.MemberCode == "000000" {
						log.Printf("skipped entry %q\n", row.RecordGUID)

						continue
					}

					member, err := findMemberByGUID(tx, row.MemberGUID)
					if err != nil {
						return err
					}

					record, err := findOrInitializeByLegacyID(tx, collection, row.RecordGUID)
					if err != nil {
						return err
					}

					rawRecord, err := json.Marshal(row)
					if err != nil {
						return fmt.Errorf("failed to marshal row: %w", err)
					}

					record.Set("legacy_guid", row.RecordGUID)
					record.Set("legacy_from", row.FromMonth)
					record.Set("legacy_to", row.ToMonth)
					record.Set("legacy_record", rawRecord)
					record.Set("member_id", member.GetId())
					record.Set("amount_in_euros", row.Amount)
					record.Set("months", row.Months)
					record.Set("issued_at", row.IssueDate)
					record.Set("contains_registration_fee", row.ContainsRegisterFee)

					if err := tx.SaveRecord(record); err != nil {
						return fmt.Errorf("failed to save record %q: %w", row.RecordGUID, err)
					}
				}

				return nil
			})
		},
	}
}

//nolint:funlen
func importDeletionsCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "deletions",
		Example:      "import deletions",
		SilenceUsage: true,
		PostRun: func(_ *cobra.Command, _ []string) {
			fmt.Println(importReport.String())
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no filename provided")
			}

			filename := args[0]

			sheet := "Sheet1"
			if len(args) == 2 {
				sheet = args[1]
			}

			file, err := excelize.OpenFile(filename)
			if err != nil {
				return fmt.Errorf("failed to open excel file: %w", err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					fmt.Println(err)
				}
			}()

			type Record struct {
				Code     string
				FullName string
			}

			records := []*Record{}

			rows, err := file.GetRows(sheet)
			if err != nil {
				return fmt.Errorf("failed to read from sheet %s: %w", sheet, err)
			}
			for i, row := range rows {
				if i == 0 {
					continue
				}

				records = append(records, &Record{
					Code:     strings.TrimSpace(row[0]),
					FullName: utils.Normalize(row[1]),
				})
			}

			if err := file.Close(); err != nil {
				return fmt.Errorf("failed to close excel file: %w", err)
			}

			ctx := &echo.DefaultContext{}
			ctx.Set("dao", app.Dao())

			for _, record := range records {
				member, err := memberQuery.FindByNo(ctx, record.Code)
				if err != nil {
					return fmt.Errorf("failed to find member %s: %w", record.Code, err)
				}

				fullName := fmt.Sprintf("%s %s", member.LastName, member.FullName)
				if record.FullName != fullName {
					fmt.Printf("Στο εξελ:\t%s\t|%s\n", record.Code, record.FullName)
					fmt.Printf("Στο μητρώο:\t%s\t|%s\n", member.MemberNo, fullName)
					if YesNoPrompt("Προχωράμε στην απενεργοποίηση;", false) {
						// skip
						continue
					}
					// disable
				}
			}

			return nil
		},
	}
}

func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
