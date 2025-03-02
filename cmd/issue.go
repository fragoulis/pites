package cmd

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/spf13/cobra"

	"github.com/fragoulis/setip_v2/internal/app/issue/model"
	"github.com/fragoulis/setip_v2/internal/app/issue/query"
	"github.com/fragoulis/setip_v2/internal/app/issue/service"
)

type ScanReport struct {
	RecordIDs       map[string]map[string]struct{}
	IssuesUpdated   int
	IssuesCreated   int
	IssuesAutofixed int
}

func (r *ScanReport) String() string {
	perCollection := ""
	for collection, recordIDs := range r.RecordIDs {
		perCollection += fmt.Sprintf("	%s: %d\n", collection, len(recordIDs))
	}

	return fmt.Sprintf(`
RecordsWithIssues:
%s
IssuesUpdated: %d
IssuesCreated: %d
IssuesAutofixed: %d
`,
		perCollection,
		r.IssuesUpdated,
		r.IssuesCreated,
		r.IssuesAutofixed,
	)
}

//nolint:gochecknoglobals
var scanReport = ScanReport{}

func NewIssueCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "issue",
		Short: "Manages open issues.",
	}

	scanReport.RecordIDs = make(map[string]map[string]struct{}, 0)

	command.AddCommand(issueScanCommand(app))

	return command
}

type ScanType struct {
	Collection string
	IssueType  string
	Field      string
	Value      any
	Importance string
}

//nolint:gochecknoglobals
var scans = []ScanType{
	{
		Collection: "members",
		IssueType:  "missing_address",
		Field:      "address_city_id",
		Value:      "",
		Importance: model.IssueImportanceHigh,
	},
	{
		Collection: "members",
		IssueType:  "missing_address_street",
		Field:      "address_street_id",
		Value:      "",
		Importance: model.IssueImportanceHigh,
	},
	{
		Collection: "members",
		IssueType:  "missing_address_street_no",
		Field:      "address_street_no",
		Value:      "",
		Importance: model.IssueImportanceHigh,
	},
	{
		Collection: "members",
		IssueType:  "missing_email",
		Field:      "email",
		Value:      "",
		Importance: model.IssueImportanceHigh,
	},
	{
		Collection: "members",
		IssueType:  "missing_phone",
		Field:      "mobile",
		Value:      "",
		Importance: model.IssueImportanceHigh,
	},
	{
		Collection: "members",
		IssueType:  "missing_id_card_number",
		Field:      "id_card_number",
		Value:      "",
		Importance: model.IssueImportanceLow,
	},
	{
		Collection: "members",
		IssueType:  "missing_insurance_no",
		Field:      "insurance_no",
		Value:      "",
		Importance: model.IssueImportanceLow,
	},
	{
		Collection: "companies",
		IssueType:  "missing_address",
		Field:      "address_city_id",
		Value:      "",
		Importance: model.IssueImportanceMedium,
	},
	{
		Collection: "companies",
		IssueType:  "missing_address_street",
		Field:      "address_street_id",
		Value:      "",
		Importance: model.IssueImportanceMedium,
	},
	{
		Collection: "companies",
		IssueType:  "missing_address_street_no",
		Field:      "address_street_no",
		Value:      "",
		Importance: model.IssueImportanceMedium,
	},
	{
		Collection: "companies",
		IssueType:  "missing_email",
		Field:      "email",
		Value:      "",
		Importance: model.IssueImportanceMedium,
	},
	{
		Collection: "companies",
		IssueType:  "missing_phone",
		Field:      "phone",
		Value:      "",
		Importance: model.IssueImportanceMedium,
	},
	{
		Collection: "companies",
		IssueType:  "missing_website",
		Field:      "website",
		Value:      "",
		Importance: model.IssueImportanceLow,
	},
	{
		Collection: "subscriptions",
		IssueType:  "unpaid_subscription_fee",
		Field:      "fee_paid",
		Value:      false,
		Importance: model.IssueImportanceMedium,
	},
}

func performScan(
	dao *daos.Dao,
	issueTypesByKey model.IssueTypesByKey,
	issues model.Issues,
	scanType ScanType,
) ([]*model.Issue, error) {
	fmt.Printf("Scanning %s for %s\n", scanType.Collection, scanType.IssueType)

	records, err := dao.FindRecordsByExpr(scanType.Collection,
		dbx.HashExp{scanType.Field: scanType.Value},
	)
	if err != nil {
		return nil, err
	}

	issueType, ok := issueTypesByKey[scanType.IssueType]
	if !ok {
		return nil, fmt.Errorf("invalid issue type: %s", scanType.IssueType)
	}

	for _, record := range records {
		issues = append(issues, &model.Issue{
			IssueTypeID:  issueType.ID,
			RelationName: scanType.Collection,
			RelationID:   record.GetId(),
			Importance:   scanType.Importance,
		})

		_, ok := scanReport.RecordIDs[scanType.Collection]
		if !ok {
			scanReport.RecordIDs[scanType.Collection] = make(map[string]struct{}, 0)
		}

		scanReport.RecordIDs[scanType.Collection][record.GetId()] = struct{}{}
	}

	return issues, nil
}

func generateIssues(dao *daos.Dao, issues model.Issues) error {
	fmt.Println("Generating issues...")

	for _, issue := range issues {
		created, err := service.CreateIssue(dao, issue)
		if err != nil {
			//nolint:wrapcheck
			return err
		}

		if created {
			scanReport.IssuesCreated++
		} else {
			scanReport.IssuesUpdated++
		}
	}

	return nil
}

func generateReport() {
	fmt.Println(scanReport.String())
}

func scanForIssues(dao *daos.Dao) error {
	issueTypesByKey, err := query.MapIssueTypesByKey(dao)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	var issues model.Issues

	for _, scan := range scans {
		issues, err = performScan(
			dao,
			issueTypesByKey,
			issues,
			scan,
		)
		if err != nil {
			return err
		}
	}

	err = generateIssues(dao, issues)
	if err != nil {
		return err
	}

	generateReport()

	return nil
}

func issueScanCommand(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "scan",
		Example:      "issue scan",
		Short:        "Scans database for issues and registers them.",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.Dao().RunInTransaction(scanForIssues)
		},
	}
}
