package model

import (
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type (
	IssueTypesByKey map[string]*IssueType
	Issues          []*Issue
)

const (
	IssueImportanceLow    = "low"
	IssueImportanceMedium = "medium"
	IssueImportanceHigh   = "high"
)

type IssueType struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

type Issue struct {
	ID           string    `json:"id"`
	IssueTypeID  string    `json:"issue_type_id"`
	RelationName string    `json:"relation_name"`
	RelationID   string    `json:"relation_id"`
	ResolvedAt   time.Time `json:"resolved_at"`
	Importance   string    `json:"importance"`
	Key          string    `json:"key"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

func (i *Issue) IsResolved() bool {
	return !i.ResolvedAt.IsZero()
}

func NewFromRecord(rec *models.Record) *Issue {
	model := &Issue{
		ID:           rec.GetId(),
		IssueTypeID:  rec.GetString("issue_type_id"),
		RelationName: rec.GetString("relation_name"),
		RelationID:   rec.GetString("relation_id"),
		ResolvedAt:   rec.GetDateTime("resolved_at").Time(),
		Importance:   rec.GetString("importance"),
		CreatedAt:    rec.GetDateTime("created").Time(),
	}

	issueType := rec.ExpandedOne("issue_type_id")
	if issueType != nil {
		model.Key = issueType.GetString("key")
		model.Description = issueType.GetString("description")
	}

	return model
}

func NewIssueTypeFromRecord(rec *models.Record) *IssueType {
	return &IssueType{
		ID:          rec.GetId(),
		Key:         rec.GetString("key"),
		Description: rec.GetString("description"),
	}
}
