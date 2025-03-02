package query

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/employment"
	"github.com/fragoulis/setip_v2/internal/app/errors"
	issueModel "github.com/fragoulis/setip_v2/internal/app/issue/model"
	issueQuery "github.com/fragoulis/setip_v2/internal/app/issue/query"
	"github.com/fragoulis/setip_v2/internal/app/member/model"
	paymentModel "github.com/fragoulis/setip_v2/internal/app/payment/model"
)

type (
	EmploymentsByMemberID   map[string][]*employment.Employment
	SubscriptionsByMemberID map[string][]*model.Subscription
	PaymentsByMemberID      map[string][]*paymentModel.Payment
	IssuesByMemberID        map[string][]*issueModel.Issue
)

func Search(ctx echo.Context, params *SearchParams) ([]*model.Member, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	query := params.Apply(dao.RecordQuery("members"))

	records := []*models.Record{}

	err := query.
		OrderBy("last_name ASC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	slog.Info("member", "action", "search", "count", len(records))

	return enrichMemberRecords(ctx, dao, records)
}

func Count(ctx echo.Context, params *SearchParams) (int, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return 0, errors.ErrFailedToGetDao
	}

	query := params.Apply(dao.RecordQuery("members"))

	var count int

	err := query.
		Select("count(distinct members.id)").
		Limit(1).Row(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return count, nil
}

func FindByID(ctx echo.Context, id string) (*model.Member, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	record, err := dao.FindRecordById("members", id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	members, err := enrichMemberRecords(ctx, dao, []*models.Record{record})
	if err != nil {
		return nil, err
	}

	return members[0], nil
}

func FindByNo(ctx echo.Context, memberNo string) (*model.Member, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records, err := dao.FindRecordsByExpr("members", dbx.HashExp{"member_no": memberNo})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	members, err := enrichMemberRecords(ctx, dao, records)
	if err != nil {
		return nil, err
	}

	return members[0], nil
}

func NextMemberNo(ctx echo.Context) (int, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return 0, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("members").
		OrderBy("member_no DESC").
		Limit(1).
		All(&records)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(records) == 0 {
		return 1, nil
	}

	memberNo := records[0].GetInt("member_no")

	return memberNo + 1, nil
}

//nolint:cyclop
func enrichMemberRecords(ctx echo.Context, dao *daos.Dao, records []*models.Record) ([]*model.Member, error) {
	// Expand the records relations (aka load associations).
	if ctx.Request() != nil {
		err := apis.EnrichRecords(
			ctx,
			dao,
			records,
			"address_city_id",
			"address_street_id",
			"company_id.parent_id.business_type_id",
			"company_id.business_type_id",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to expand relations for members: %w", err)
		}
	} else {
		errs := dao.ExpandRecords(records, []string{"address_city_id", "address_street_id"}, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to expand relations for members: %s", errs)
		}
	}

	memberIDs := []string{}
	for _, record := range records {
		memberIDs = append(memberIDs, record.GetId())
	}

	employmentsByMemberID, err := getEmploymentsByMemberID(ctx, memberIDs)
	if err != nil {
		return nil, err
	}

	subscriptionsByMemberID, err := getSubscriptionsByMemberID(ctx, memberIDs)
	if err != nil {
		return nil, err
	}

	paymentsByMemberID, err := getPaymentsByMemberID(ctx, memberIDs)
	if err != nil {
		return nil, err
	}

	issuesByMemberID, err := getIssuesByMemberID(ctx, memberIDs)
	if err != nil {
		return nil, err
	}

	members := []*model.Member{}

	for _, record := range records {
		employments, ok := employmentsByMemberID[record.GetId()]
		if !ok {
			employments = []*employment.Employment{}
		}

		subscriptions, ok := subscriptionsByMemberID[record.GetId()]
		if !ok {
			subscriptions = []*model.Subscription{}
		}

		payments, ok := paymentsByMemberID[record.GetId()]
		if !ok {
			payments = []*paymentModel.Payment{}
		}

		issues, ok := issuesByMemberID[record.GetId()]
		if !ok {
			issues = []*issueModel.Issue{}
		}

		members = append(members, model.NewFromRecord(
			record,
			employments,
			subscriptions,
			payments,
			issues,
		))
	}

	return members, nil
}

func getEmploymentsByMemberID(ctx echo.Context, ids []string) (EmploymentsByMemberID, error) {
	rels, err := employment.FindByMemberID(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find rels: %w", err)
	}

	relsByMemberID := EmploymentsByMemberID{}
	for _, rel := range rels {
		relsByMemberID[rel.MemberID] = append(
			relsByMemberID[rel.MemberID],
			rel,
		)
	}

	return relsByMemberID, nil
}

func getSubscriptionsByMemberID(ctx echo.Context, ids []string) (SubscriptionsByMemberID, error) {
	rels, err := findSubscriptionsByMemberID(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find rels: %w", err)
	}

	relsByMemberID := SubscriptionsByMemberID{}
	for _, rel := range rels {
		relsByMemberID[rel.MemberID] = append(
			relsByMemberID[rel.MemberID],
			rel,
		)
	}

	return relsByMemberID, nil
}

func getPaymentsByMemberID(ctx echo.Context, ids []string) (PaymentsByMemberID, error) {
	rels, err := paymentModel.FindByMemberID(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find rels: %w", err)
	}

	relsByMemberID := PaymentsByMemberID{}
	for _, rel := range rels {
		relsByMemberID[rel.MemberID] = append(
			relsByMemberID[rel.MemberID],
			rel,
		)
	}

	return relsByMemberID, nil
}

func getIssuesByMemberID(ctx echo.Context, ids []string) (IssuesByMemberID, error) {
	rels, err := issueQuery.FindUnresolvedByRelationID(ctx, "members", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find rels: %w", err)
	}

	relsByMemberID := IssuesByMemberID{}
	for _, rel := range rels {
		relsByMemberID[rel.RelationID] = append(
			relsByMemberID[rel.RelationID],
			rel,
		)
	}

	return relsByMemberID, nil
}

func findSubscriptionsByMemberID(ctx echo.Context, memberIDs []string) ([]*model.Subscription, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("subscriptions").
		Where(dbx.In("member_id", list.ToInterfaceSlice(memberIDs)...)).
		OrderBy("start_date DESC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	models := []*model.Subscription{}
	for _, record := range records {
		models = append(models, model.NewSubscriptionFromRecord(record))
	}

	return models, nil
}
