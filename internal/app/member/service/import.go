package service

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/xuri/excelize/v2"

	companymodel "github.com/fragoulis/setip_v2/internal/app/company/model"
	companyquery "github.com/fragoulis/setip_v2/internal/app/company/query"
	companyservice "github.com/fragoulis/setip_v2/internal/app/company/service"
	membermodel "github.com/fragoulis/setip_v2/internal/app/member/model"
	memberquery "github.com/fragoulis/setip_v2/internal/app/member/query"
	paymentservice "github.com/fragoulis/setip_v2/internal/app/payment/service"
	dbaddress "github.com/fragoulis/setip_v2/internal/db/address"
	"github.com/fragoulis/setip_v2/internal/utils"
)

//nolint:gochecknoglobals
var (
	headers = []string{
		"Α/Μ",
		"Όνομα",
		"Επώνυμο",
		"Πατρώνυμο",
		"Δήμος",
		"Οδός",
		"Αριθμός",
		"Ομιλος",
		"Εταιρεία",
		"Δημός Εταιρείας",
		"Οδός Εταιρείας",
		"Αριθμός Εταιρείας",
		"Τυπος",
		"Ειδικότητα",
		"Email",
		"Κινητό",
		"Σταθερό",
		"Ημ/νία γέννησης",
		"ΑΔΤ",
		"ΑΜΑ",
		"Σχόλια",
		"Ημ/νία εγγραφής",
		"Έχει πληρώσει μέχρι",
	}

	companyGroupIndex    = 7
	companyNameIndex     = 8
	companyCityIndex     = 9
	companyStreetIndex   = 10
	companyStreetNoIndex = 11
)

type (
	CompanyGroupsByName map[string]*companymodel.Company
	CompaniesByName     map[string]*companymodel.Company
	MembersByNo         map[int]*membermodel.Member
)

func Import(ctx echo.Context, app *pocketbase.PocketBase, src io.Reader) (MembersByNo, CompaniesByName, error) {
	file, err := excelize.OpenReader(src)
	if err != nil {
		return nil, nil, fmt.Errorf("read excel: %w", err)
	}

	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			app.Logger().Error("close excel", "error", err.Error())
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, nil, fmt.Errorf("get rows: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil, errors.New("empty spreadsheet")
	}

	if !areHeadersValid(app, rows[0]) {
		return nil, nil, errors.New("invalid headers")
	}

	existingCompanies, err := companyquery.Search(ctx, app.Dao(), &companyquery.ListCompaniesRequest{
		Limit: 1000,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("load copmanies: %w", err)
	}

	existingCompaniesByName := make(CompaniesByName)

	for _, company := range existingCompanies {
		key := company.Name
		if company.ParentID != "" {
			key = company.Name + company.Branch
		}

		existingCompaniesByName[key] = company
	}

	existingMembers, err := memberquery.Search(ctx, &memberquery.SearchParams{
		Limit: 10000,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("load copmanies: %w", err)
	}

	existingMembersByNo := make(MembersByNo)

	for _, member := range existingMembers {
		no, _ := strconv.Atoi(member.MemberNo)
		existingMembersByNo[no] = member
	}

	var membersByNo MembersByNo
	var companiesByName CompaniesByName

	// Persist data
	err = app.Dao().RunInTransaction(func(tx *daos.Dao) error {
		ctx.Set("dao", tx)

		companyGroupsByName, err := persistCompanyGroups(ctx, app, tx, rows, existingCompaniesByName)
		if err != nil {
			return fmt.Errorf("persist company groups: %w", err)
		}

		companiesByName, err = persistCompanies(ctx, app, tx, rows, companyGroupsByName, existingCompaniesByName)
		if err != nil {
			return fmt.Errorf("persist companies: %w", err)
		}

		membersByNo, err = persistMembers(ctx, app, tx, rows, companiesByName, existingMembersByNo)
		if err != nil {
			return fmt.Errorf("persist members: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("persist data: %w", err)
	}

	return membersByNo, companiesByName, nil
}

func areHeadersValid(app *pocketbase.PocketBase, row []string) bool {
	for idx, cell := range row {
		if utils.Equal(headers[idx], cell) {
			continue
		}

		app.Logger().Error("headers do not match",
			"idx", idx,
			"src", utils.Normalize(headers[idx]),
			"in", utils.Normalize(cell),
			"row", row,
		)

		return false
	}

	return true
}

func persistCompanyGroups(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	tx *daos.Dao,
	rows [][]string,
	existingCompaniesByName CompaniesByName,
) (CompanyGroupsByName, error) {
	companyGroups := make(map[string]struct{})

	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		companyGroups[utils.Normalize(row[companyGroupIndex])] = struct{}{}
	}

	companyGroupsByName := make(CompanyGroupsByName)

	for name := range companyGroups {
		if companyGroupRecord, exists := existingCompaniesByName[name]; exists {
			companyGroupsByName[companyGroupRecord.Name] = companyGroupRecord

			continue
		}

		companyGroupRecord, err := companyservice.Create(ctx, app, tx, map[string]any{
			"name": name,
		})
		if err != nil {
			return nil, fmt.Errorf("create company group: %w", err)
		}

		companyGroupsByName[companyGroupRecord.Name] = companyGroupRecord
	}

	return companyGroupsByName, nil
}

func persistCompanies(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	tx *daos.Dao,
	rows [][]string,
	companyGroupsByName CompanyGroupsByName,
	existingCompaniesByName CompaniesByName,
) (CompaniesByName, error) {
	type companySpec struct {
		cityID   string
		streetID string
		streetNo string
		parentID string
		parent   string
	}

	companies := make(map[string]companySpec)

	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		groupName := utils.Normalize(row[companyGroupIndex])
		parentID := companyGroupsByName[groupName].ID

		cityName := utils.Normalize(row[companyCityIndex])

		city, err := dbaddress.FindCityByName(tx, cityName)
		if err != nil {
			return nil, fmt.Errorf("find city: %w", err)
		}

		streetID := ""
		streetName := utils.Normalize(row[companyStreetIndex])

		if streetName != "" {
			street, err := dbaddress.FindStreetByName(tx, streetName, city)
			if err != nil {
				return nil, fmt.Errorf("find street: %w", err)
			}

			streetID = street.GetId()
		}

		companies[utils.Normalize(row[companyNameIndex])] = companySpec{
			parent:   groupName,
			parentID: parentID,
			cityID:   city.GetId(),
			streetID: streetID,
			streetNo: utils.Normalize(row[companyStreetNoIndex]),
		}
	}

	companyRecords := make(CompaniesByName)

	for name, spec := range companies {
		key := name
		if spec.parentID != "" {
			key = spec.parent + name
		}

		if companyRecord, exists := existingCompaniesByName[key]; exists {
			companyRecords[key] = companyRecord

			continue
		}

		companyRecord, err := companyservice.Create(ctx, app, tx, map[string]any{
			"name":              name,
			"parent_id":         spec.parentID,
			"address_street_id": spec.streetID,
			"address_city_id":   spec.cityID,
			"address_street_no": spec.streetNo,
		})
		if err != nil {
			return nil, fmt.Errorf("create company %q: %w", key, err)
		}

		companyRecords[companyRecord.Name] = companyRecord
	}

	return companyRecords, nil
}

func persistMembers(
	ctx echo.Context,
	app *pocketbase.PocketBase,
	tx *daos.Dao,
	rows [][]string,
	companiesByName CompaniesByName,
	existingMembersByNo MembersByNo,
) (MembersByNo, error) {
	membersByNo := make(MembersByNo)

	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		// Company
		companyName := utils.Normalize(row[companyNameIndex])
		companyGroup := utils.Normalize(row[companyGroupIndex])

		companyKey := companyName
		if companyGroup != "" {
			companyKey = companyGroup + companyName
		}

		company, ok := companiesByName[companyKey]
		if !ok {
			return nil, fmt.Errorf("company %q not found", companyKey)
		}

		companyID := company.ID

		// City
		cityName := utils.Normalize(row[4])

		city, err := dbaddress.FindCityByName(tx, cityName)
		if err != nil {
			return nil, fmt.Errorf("find city: %w", err)
		}

		// Street
		streetID := ""

		streetName := utils.Normalize(row[5])
		if streetName != "" {
			street, err := dbaddress.FindStreetByName(tx, streetName, city)
			if err != nil {
				return nil, fmt.Errorf("find street: %w", err)
			}

			streetID = street.GetId()
		}

		// Birthdate
		birthdate, err := time.Parse("02/01/2006", utils.Normalize(row[17]))
		if err != nil {
			return nil, fmt.Errorf("parse birthdate: %w", err)
		}

		// Registration
		registration, err := time.Parse("01/2006", utils.Normalize(row[21]))
		if err != nil {
			return nil, fmt.Errorf("parse registration: %w", err)
		}

		// Paid until
		paidUntil, err := time.Parse("01/2006", utils.Normalize(row[22]))
		if err != nil {
			return nil, fmt.Errorf("parse paidUntil: %w", err)
		}

		memberNo, err := strconv.Atoi(utils.Normalize(row[0]))
		if err != nil {
			return nil, fmt.Errorf("member no %q is not a number", row[0])
		}

		updateRequest := UpdateMemberRequest{
			FirstName:         utils.Normalize(row[1]),
			LastName:          utils.Normalize(row[2]),
			FatherName:        utils.Normalize(row[3]),
			AddressCityID:     city.GetId(),
			AddressStreetID:   streetID,
			AddressStreetNo:   utils.Normalize(row[6]),
			CompanyID:         companyID,
			Specialty:         utils.Normalize(row[13]),
			Email:             utils.Normalize(row[14]),
			Mobile:            utils.Normalize(row[15]),
			Phone:             utils.Normalize(row[16]),
			Birthdate:         birthdate.Format("2006-01-02"),
			IDCardNumber:      utils.Normalize(row[18]),
			SocialSecurityNum: utils.Normalize(row[19]),
			Comments:          row[20],
		}

		var rec *membermodel.Member

		if existingMember, exists := existingMembersByNo[memberNo]; exists {
			existingRecord, err := tx.FindRecordById("members", existingMember.ID)
			if err != nil {
				return nil, fmt.Errorf("find member: %w", err)
			}

			rec, err = UpdateMember(ctx, app, tx, existingRecord, &updateRequest)
			if err != nil {
				return nil, fmt.Errorf("update member: %w", err)
			}
		} else {
			createRequest := &CreateMemberRequest{
				MemberNo:            memberNo,
				UpdateMemberRequest: updateRequest,
				CreateSubscriptionRequest: CreateSubscriptionRequest{
					StartDate: registration.Format("2006-01-02"),
					FeePaid:   true,
				},
			}

			rec, err = CreateMember(ctx, app, tx, createRequest)
			if err != nil {
				return nil, fmt.Errorf("create member: %w", err)
			}
		}

		// Delete present payments for member
		_, err = tx.DB().
			NewQuery(fmt.Sprintf("DELETE FROM payments WHERE member_id = '%s'", rec.ID)).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("delete payments: %w", err)
		}

		months, err := utils.MonthsSince(registration, paidUntil.AddDate(0, 1, 0))
		if err != nil {
			return nil, fmt.Errorf("invalid paid until: %w", err)
		}

		_, err = paymentservice.Create(ctx, app, &paymentservice.CreatePaymentRequest{
			MemberID: rec.ID,
			Months:   months,
			IssuedAt: time.Now().Format(time.DateOnly),
			Comments: "ΑΡΧΙΚΗ ΕΙΣΑΓΩΓΗ",
		})
		if err != nil {
			return nil, fmt.Errorf("create payment %v: %w", row, err)
		}

		membersByNo[memberNo] = rec
	}

	return membersByNo, nil
}
