package service

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"

	"github.com/fragoulis/setip_v2/internal/app/member/model"
)

//nolint:funlen,cyclop,mnd
func Export(destFile *os.File, members []*model.Member, columns []string) error {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	evenStyleID, err := file.NewStyle(createRowStyle("f7f7f7"))
	if err != nil {
		return fmt.Errorf("failed to create even style: %w", err)
	}

	oddStyleID, err := file.NewStyle(createRowStyle("ffffff"))
	if err != nil {
		return fmt.Errorf("failed to create odd style: %w", err)
	}

	sheet := "Sheet1"

	// Set column width depending on the approximate data length.
	// The following code will generate columns such as:
	//  _ = file.SetColWidth(sheet, "A", "B", 10)
	//	_ = file.SetColWidth(sheet, "B", "C", 40)
	//	_ = file.SetColWidth(sheet, "C", "D", 40)
	//	... etc
	startCol := 'A'

	for idx := range columns {
		width := 40
		if idx == 0 {
			width = 10
		}

		_ = file.SetColWidth(sheet, string(startCol+rune(idx)), string(startCol+rune(idx+1)), float64(width))
	}

	// Create the header row
	headers := []string{}

	// Break the address column into 3 granular columns for city, street and number.
	for _, column := range columns {
		switch column {
		case "Διεύθυνση":
			headers = append(headers, "Δήμος")
			headers = append(headers, "Οδός/Αριθμός")
			headers = append(headers, "ΤΚ")
		default:
			headers = append(headers, column)
		}
	}

	if err := file.SetSheetRow(sheet, "A1", &headers); err != nil {
		return fmt.Errorf("failed to set sheet headers: %w", err)
	}

	// Set the height for header row
	if err := file.SetRowHeight(sheet, 1, 20); err != nil {
		return fmt.Errorf("failed to set row height: %w", err)
	}

	// Create the member rows
	for idx, member := range members {
		row := idx + 2

		mem := []string{}

		for _, column := range columns {
			switch column {
			case "Α/Μ":
				mem = append(mem, member.MemberNo)
			case "Όνομα":
				mem = append(mem, fmt.Sprintf("%s %s του %s", member.LastName, member.FirstName, member.FatherName))
			case "Διεύθυνση":
				mem = append(mem, cityWithFallback(member))
				mem = append(mem, streetWithFallback(member))
				mem = append(mem, postCodeWithFallback(member))
			case "Κινητό":
				mem = append(mem, member.Mobile)
			case "Email":
				mem = append(mem, member.Email)
			case "Συνδρομή":
				mem = append(mem, member.SubscriptionFormatted)
			case "Οικονομικά":
				mem = append(mem, member.PaymentStatus.Formatted)
			case "Ομάδα":
				mem = append(mem, member.BusinessTypeName)
			case "Εταιρεία":
				mem = append(mem, member.CompanyName)
			case "Παράρτημα":
				mem = append(mem, member.CompanyBranchName)
			case "Δ/ση Εταιρείας":
				mem = append(mem, member.CompanyAddress)
			case "ΑΔΤ":
				mem = append(mem, member.IDCardNumber)
			}
		}

		if err := file.SetSheetRow(sheet, fmt.Sprintf("%s%d", "A", row), &mem); err != nil {
			return fmt.Errorf("failed to set sheet row %d: %w", row, err)
		}

		if err := file.SetRowHeight(sheet, row, 20); err != nil {
			return fmt.Errorf("failed to set row %d height: %w", row, err)
		}

		if idx%2 == 0 {
			if err := file.SetRowStyle(sheet, row, row, evenStyleID); err != nil {
				return fmt.Errorf("failed to set row %d style: %w", row, err)
			}

			continue
		}

		if err := file.SetRowStyle(sheet, row, row, oddStyleID); err != nil {
			return fmt.Errorf("failed to set row %d style: %w", row, err)
		}
	}

	// Save spreadsheet by the given path.
	if err := file.SaveAs(destFile.Name()); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func createRowStyle(fillColor string) *excelize.Style {
	return &excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 0},
			{Type: "top", Color: "e0e0e0", Style: 1},
			{Type: "bottom", Color: "e0e0e0", Style: 1},
			{Type: "right", Color: "FF0000", Style: 0},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{fillColor},
			Pattern: 1,
		},
	}
}

func streetWithFallback(member *model.Member) string {
	if member.AddressStreetName != "" {
		return fmt.Sprintf("%s %s", member.AddressStreetName, member.AddressStreetNo)
	}

	return member.LegacyAddress
}

func cityWithFallback(member *model.Member) string {
	if member.AddressCityName != "" {
		return member.AddressCityName
	}

	if member.LegacyCity != "" {
		return member.LegacyCity
	}

	return member.LegacyArea
}

func postCodeWithFallback(member *model.Member) string {
	if member.AddressPostCode != "" {
		return member.AddressPostCode
	}

	return member.LegacyPostCode
}
