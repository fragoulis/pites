package model

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/errors"
	"github.com/fragoulis/setip_v2/internal/utils"
)

type Payment struct {
	ID                string    `json:"id"`
	MemberID          string    `json:"member_id"`
	MemberNo          string    `json:"member_no"`
	MemberName        string    `json:"member_name"`
	Amount            int       `json:"amount"`
	Months            int       `json:"months"`
	ReceiptBlockNo    int       `json:"receipt_block_no"`
	ReceiptNo         int       `json:"receipt_no"`
	CreatedByUserID   string    `json:"created_by_user_id"`
	CreatedByUser     *User     `json:"created_by_user"`
	IssuedAt          time.Time `json:"issued_at"`
	IssuedAtFormatted string    `json:"issued_at_formatted"`
	Comments          string    `json:"comments"`
	LegacyTo          time.Time `json:"legacy_to"`
	LegacyToFormatted string    `json:"legacy_to_formatted"`
}

func NewFromRecord(rec *models.Record) *Payment {
	member := rec.ExpandedOne("member_id")

	return NewFromRecordNoMember(
		rec,
		fmt.Sprintf("%06d", member.GetInt("member_no")),
		member.GetString("full_name"),
	)
}

func NewFromRecordNoMember(rec *models.Record, memberNo, memberName string) *Payment {
	return &Payment{
		ID:                rec.GetId(),
		MemberID:          rec.GetString("member_id"),
		MemberNo:          memberNo,
		MemberName:        memberName,
		Amount:            rec.GetInt("amount_in_euros"),
		Months:            rec.GetInt("months"),
		ReceiptBlockNo:    rec.GetInt("receipt_block_no"),
		ReceiptNo:         rec.GetInt("receipt_no"),
		CreatedByUserID:   rec.GetString("created_by_user_id"),
		CreatedByUser:     newUserFromRecord(rec.ExpandedOne("created_by_user_id")),
		IssuedAt:          rec.GetDateTime("issued_at").Time(),
		IssuedAtFormatted: utils.Day(rec.GetDateTime("issued_at").Time()),
		Comments:          rec.GetString("comments"),
		LegacyTo:          rec.GetDateTime("legacy_to").Time(),
		LegacyToFormatted: utils.Month(rec.GetDateTime("legacy_to").Time()),
	}
}

func FindByMemberID(ctx echo.Context, memberIDs []string) ([]*Payment, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("payments").
		Where(dbx.In("member_id", list.ToInterfaceSlice(memberIDs)...)).
		OrderBy("issued_at DESC").
		All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Expand the records relations (aka load associations).
	if ctx.Request() != nil {
		err = apis.EnrichRecords(
			ctx,
			dao,
			records,
			"member_id",
			"created_by_user_id",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to expand relations for payments: %w", err)
		}
	} else {
		errs := dao.ExpandRecords(records, []string{"created_by_user_id"}, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to expand relations for payments: %v", errs)
		}
	}

	models := make([]*Payment, 0, len(records))
	for _, record := range records {
		models = append(models, NewFromRecord(record))
	}

	return models, nil
}
