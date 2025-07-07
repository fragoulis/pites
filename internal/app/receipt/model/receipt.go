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

type Receipt struct {
	ID                string    `json:"id"`
	MemberID          string    `json:"member_id"`
	MemberNo          string    `json:"member_no"`
	MemberName        string    `json:"member_name"`
	Amount            int       `json:"amount"`
	BlockNo           int       `json:"block_no"`
	ReceiptNo         int       `json:"receipt_no"`
	CreatedByUserID   string    `json:"created_by_user_id"`
	CreatedByUser     *User     `json:"created_by_user"`
	IssuedAt          time.Time `json:"issued_at"`
	IssuedAtFormatted string    `json:"issued_at_formatted"`
	Comments          string    `json:"comments"`
}

func NewFromRecord(rec *models.Record) *Receipt {
	member := rec.ExpandedOne("member_id")

	return NewFromRecordNoMember(
		rec,
		fmt.Sprintf("%06d", member.GetInt("member_no")),
		member.GetString("full_name"),
	)
}

func NewFromRecordNoMember(rec *models.Record, memberNo, memberName string) *Receipt {
	return &Receipt{
		ID:                rec.GetId(),
		MemberID:          rec.GetString("member_id"),
		MemberNo:          memberNo,
		MemberName:        memberName,
		Amount:            rec.GetInt("amount_in_euros"),
		BlockNo:           rec.GetInt("block_no"),
		ReceiptNo:         rec.GetInt("receipt_no"),
		CreatedByUserID:   rec.GetString("created_by_user_id"),
		CreatedByUser:     newUserFromRecord(rec.ExpandedOne("created_by_user_id")),
		IssuedAt:          rec.GetDateTime("issued_at").Time(),
		IssuedAtFormatted: utils.Day(rec.GetDateTime("issued_at").Time()),
		Comments:          rec.GetString("comments"),
	}
}

func FindByMemberID(ctx echo.Context, memberIDs []string) ([]*Receipt, error) {
	dao, ok := ctx.Get("dao").(*daos.Dao)
	if !ok {
		return nil, errors.ErrFailedToGetDao
	}

	records := []*models.Record{}

	err := dao.RecordQuery("receipts").
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
			return nil, fmt.Errorf("failed to expand relations for receipts: %w", err)
		}
	} else {
		errs := dao.ExpandRecords(records, []string{"created_by_user_id"}, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to expand relations for receipts: %v", errs)
		}
	}

	models := make([]*Receipt, 0, len(records))
	for _, record := range records {
		models = append(models, NewFromRecord(record))
	}

	return models, nil
}
