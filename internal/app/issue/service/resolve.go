package service

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"

	"github.com/fragoulis/setip_v2/internal/app/issue/query"
)

func ResolveIssuesForMember(dao *daos.Dao, memberID string, keys ...string) error {
	issueTypeIDs, err := query.FindIssueTypeIDsByKey(
		dao,
		keys...,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	var records []*models.Record

	err = dao.RecordQuery("issues").
		Where(dbx.HashExp{"member_id": memberID}).
		Where(dbx.HashExp{"resolved_at": ""}).
		Where(dbx.In("issue_type_id", list.ToInterfaceSlice(issueTypeIDs)...)).
		All(&records)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return dao.RunInTransaction(func(tx *daos.Dao) error {
		for _, record := range records {
			record.Set("resolved_at", time.Now())

			err := tx.SaveRecord(record)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func ResolveAddressIssuesForMember(dao *daos.Dao, memberID string) error {
	return ResolveIssuesForMember(
		dao,
		memberID,
		"missing_address",
		"missing_address_street",
		"missing_address_street_no",
		"needs_confirmation_address",
	)
}
