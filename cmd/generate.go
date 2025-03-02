//nolint
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fragoulis/setip_v2/internal/utils"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

type ExcelRow struct {
	Code          string
	FullName      string
	FatherName    string
	Company       string
	RegisterDate  string
	Mobile        string
	Phone         string
	Education     string
	Specialty     string
	CompanyBranch string
	Birthdate     string
	Address       string
	Area          string
	City          string
	EMail         string
}

type MemberDataCSVRow struct {
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
	RecordGuid        string `json:"RecordGuid"`
	Phone             string `json:"Phone"`
	Mobile            string `json:"Mobile"`
	EMail             string `json:"EMail"`
	CompanyGuid       string `json:"CompanyGuid"`
	CompanyName       string `json:"CompanyName"`
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

func atoi(s string) int {
	if s == "" {
		return 0
	}

	o, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("atoi failed %q: %w", s, err))
	}

	return o
}

func atob(s string) bool {
	return s == "1"
}

func nullToEmpty(s string) string {
	if s == "NULL" {
		return ""
	}

	return s
}

func NewGenerateCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "generate",
		Short: "Generates jsom from other formats",
	}

	command.AddCommand(generateJsonMembersFromExcel(app))
	command.AddCommand(generateJsonCompaniesFromExcel(app))
	command.AddCommand(generateJsonMembersFromCSV(app))

	return command
}

// Generate json from excel.
func generateJsonMembersFromExcel(_ core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "membersfromexcel",
		Example:      "generate membersfromexcel",
		Short:        "Generate json files from excel file.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]
			dest := args[1]

			// First read all rows from the file.
			excelRows, err := readRowsFromExcel(filename)
			if err != nil {
				return err
			}

			return generateMemberJSONFromExcel(excelRows, dest)
		},
	}
}

// Generate json from excel.
func generateJsonCompaniesFromExcel(_ core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "companiesfromexcel",
		Example:      "generate companiesfromexcel",
		Short:        "Generate json files from excel file.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]
			dest := args[1]

			// First read all rows from the file.
			excelRows, err := readRowsFromExcel(filename)
			if err != nil {
				return err
			}

			return generateCompanyJSONFromExcel(excelRows, dest)
		},
	}
}

func readRowsFromExcel(filename string) ([]*ExcelRow, error) {
	var excelRows []*ExcelRow

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		member := &ExcelRow{
			Code:         row[0],
			FullName:     row[1],
			FatherName:   row[2],
			Company:      row[3],
			RegisterDate: row[4],
			Mobile:       strings.Trim(row[5], "'"),
		}

		if len(row) > 6 {
			member.Phone = row[6]
			member.Education = row[7]
			member.Specialty = row[8]
			member.CompanyBranch = row[9]
		}

		if len(row) > 10 {
			member.Birthdate = row[10]
		}

		if len(row) > 13 {
			member.Address = row[13]
		}

		if len(row) > 14 {
			member.Area = row[14]
		}

		if len(row) > 15 {
			member.City = row[15]
		}

		if len(row) > 16 {
			member.EMail = row[16]
		}

		member.Company = utils.Normalize(member.Company)

		if member.Company == "OTE" {
			member.Company = "ΟΤΕ" // to greek
		}

		member.CompanyBranch = utils.Normalize(member.CompanyBranch)
		if member.CompanyBranch == utils.Normalize("ΑΓΝΩΣΤΟ") {
			member.CompanyBranch = ""
		}

		excelRows = append(excelRows, member)
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close excel file: %w", err)
	}

	fmt.Printf("Generating %d members\n", len(excelRows))

	return excelRows, nil
}

func generateCompanyJSONFromExcel(excelRows []*ExcelRow, dest string) error {
	companiesByName := make(map[string]*CompanyJSONRow)

	for _, row := range excelRows {
		if row.Company == utils.Normalize("ΑΝΕΡΓΟΣ") {
			continue
		}

		key := row.Company + row.CompanyBranch

		_, ok := companiesByName[key]
		if ok {
			continue
		}

		company := &CompanyJSONRow{
			FullName:    row.Company,
			MainAddress: "",
			Branch:      row.CompanyBranch,
			Phone1:      "",
			Email:       "",
			WebSite:     "",
			IsActive:    true,
			RecordGUID:  "",
		}

		companiesByName[key] = company
	}

	companies := []*CompanyJSONRow{}
	for _, company := range companiesByName {
		companies = append(companies, company)
	}

	jsonData, err := json.MarshalIndent(companies, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal companies to json: %w", err)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create companies json file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write json data: %w", err)
	}

	return nil
}

func generateMemberJSONFromExcel(excelRows []*ExcelRow, dest string) error {
	members := []*MemberJSONRow{}

	for i, row := range excelRows {
		registerDate := ""
		birthDateYear := 0
		birthDate := ""

		if row.RegisterDate != "" {
			registerDateAsTime, err := time.Parse("1/2/2006", row.RegisterDate)
			if err != nil {
				registerDateAsTime, err = time.Parse("1/2/06", row.RegisterDate)
				if err != nil {
					return fmt.Errorf("failed to parse registerDate for (%d) %s: %w", i, row.FullName, err)
				}
			}

			registerDate = registerDateAsTime.Format("2006-01-02T00:00:00")
		}

		if row.Birthdate != "" {
			birthDateAsTime, err := time.Parse("1/2/2006", row.Birthdate)
			if err != nil {
				birthDateAsTime, err = time.Parse("1/2/06", row.Birthdate)
				if err != nil {
					return fmt.Errorf("failed to parse birthdate for (%d) %s: %w", i, row.FullName, err)
				}
			}

			birthDateYear = birthDateAsTime.Year()
			birthDate = birthDateAsTime.Format("2006-01-02T00:00:00")
		}

		fullName := utils.Normalize(row.FullName)
		nameParts := strings.Split(fullName, " ")

		// firstName := nameParts[0]
		// lastName := strings.Join(nameParts[1:], " ")
		firstName := strings.Join(nameParts[1:], " ")
		lastName := nameParts[0]
		fullName = fmt.Sprintf("%s %s", firstName, lastName)

		member := &MemberJSONRow{
			Code:          row.Code,
			FirstName:     utils.Normalize(firstName),
			LastName:      utils.Normalize(lastName),
			FullName:      fullName,
			FatherName:    utils.Normalize(row.FatherName),
			RegisterDate:  registerDate,
			BirthYear:     birthDateYear,
			BirthDate:     birthDate,
			Unemployed:    false,
			IsActive:      true,
			RecordGUID:    "",
			Phone:         normalizePhoneNumber(row.Phone),
			Mobile:        normalizePhoneNumber(row.Mobile),
			Email:         utils.Normalize(row.EMail),
			CompanyGUID:   "",
			CompanyName:   row.Company,
			CompanyBranch: row.CompanyBranch,
			Address:       utils.Normalize(row.Address),
			Area:          utils.Normalize(row.Area),
			City:          utils.Normalize(row.City),
			Specialty:     utils.Normalize(row.Specialty),
			Education:     utils.Normalize(row.Education),
		}

		members = append(members, member)
	}

	jsonData, err := json.MarshalIndent(members, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal members to json: %w", err)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create members json file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write json data: %w", err)
	}

	return nil
}

// Generate json from members csv.
func generateJsonMembersFromCSV(_ core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "membersfromcsv",
		Example:      "generate membersfromcsv",
		Short:        "Generate json members file from csv.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			filename := args[0]
			dest := args[1]

			// jsonRows := []MemberDataCSVRow{}

			// err := readFile("data/Members_20240928_0106.json", &jsonRows)
			// if err != nil {
			// 	return err
			// }

			csvfile, err := os.Open(filename)
			if err != nil {
				//nolint:wrapcheck
				return err
			}
			defer csvfile.Close()

			csvReader := csv.NewReader(csvfile)
			csvReader.Comma = ';'

			csvRows := []MemberDataCSVRow{}

			i := 0
			for {
				row, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("failed to read csv: %w", err)
				}

				if i == 0 {
					i++
					continue
				}

				csvRow := MemberDataCSVRow{
					Code:              row[1],
					FirstName:         utils.Normalize(nullToEmpty(row[2])),
					LastName:          utils.Normalize(nullToEmpty(row[3])),
					FullName:          utils.Normalize(nullToEmpty(row[4])),
					FatherName:        utils.Normalize(nullToEmpty(row[5])),
					RegisterDate:      nullToEmpty(row[6]),
					BirthYear:         atoi(nullToEmpty(row[7])),
					ValidThrough:      nullToEmpty(row[8]),
					Unemployed:        atob(nullToEmpty(row[9])),
					IsActive:          atob(nullToEmpty(row[10])),
					RecordGuid:        strings.ToLower(nullToEmpty(row[14])),
					Address:           utils.Normalize(nullToEmpty(row[15])),
					PostCode:          utils.Normalize(nullToEmpty(row[16])),
					Area:              utils.Normalize(nullToEmpty(row[17])),
					City:              utils.Normalize(nullToEmpty(row[18])),
					Phone:             normalizePhoneNumber(nullToEmpty(row[19])),
					Mobile:            normalizePhoneNumber(nullToEmpty(row[20])),
					EMail:             utils.Normalize(nullToEmpty(row[21])),
					IDCardNumber:      utils.Normalize(nullToEmpty(row[22])),
					IDCardAuthority:   utils.Normalize(nullToEmpty(row[24])),
					OtherUnion:        atob(nullToEmpty(row[25])),
					Specialty:         utils.Normalize(nullToEmpty(row[26])),
					Education:         utils.Normalize(nullToEmpty(row[27])),
					SocialSecurityNum: normalizePhoneNumber(nullToEmpty(row[29])),
				}

				if csvRow.IDCardAuthority == "A.T." {
					csvRow.IDCardAuthority = ""
				}

				if csvRow.FatherName != "ΆΓΝΩΣΤΟ" {
					csvRow.FatherName = ""
				}

				csvRows = append(csvRows, csvRow)
				i++
				// break
			}

			// csvRowsByCode := make(map[string]MemberDataCSVRow)
			// for _, row := range csvRows {
			// 	csvRowsByCode[row.Code] = row
			// }

			// finalJsonRows := []MemberDataCSVRow{}

			// for _, row := range jsonRows {
			// 	fromCsv := csvRowsByCode[row.Code]

			// 	row.ValidThrough = fromCsv.ValidThrough
			// 	row.Address = fromCsv.Address
			// 	row.PostCode = fromCsv.PostCode
			// 	row.Area = fromCsv.Area
			// 	row.City = fromCsv.City
			// 	row.IDCardNumber = fromCsv.IDCardNumber
			// 	row.IDCardAuthority = fromCsv.IDCardAuthority
			// 	row.OtherUnion = fromCsv.OtherUnion
			// 	row.Specialty = fromCsv.Specialty
			// 	row.Education = fromCsv.Education
			// 	row.SocialSecurityNum = fromCsv.SocialSecurityNum

			// 	finalJsonRows = append(finalJsonRows, row)
			// }

			jsonData, err := json.MarshalIndent(csvRows, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal members to json: %w", err)
			}

			file, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("failed to create members json file: %w", err)
			}
			defer file.Close()

			_, err = file.Write(jsonData)
			if err != nil {
				return fmt.Errorf("failed to write json data: %w", err)
			}

			return nil
		},
	}
}
