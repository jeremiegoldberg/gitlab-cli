package utils

import (
	"github.com/xanzy/go-gitlab"
)

// CreateMockMR creates a mock merge request for testing
func CreateMockMR(iid int, title string, description string) *gitlab.MergeRequest {
	return &gitlab.MergeRequest{
		IID:         iid,
		Title:       title,
		Description: description,
	}
}

// CreateMockIssue creates a mock issue for testing
func CreateMockIssue(iid int, title string, description string) *gitlab.Issue {
	return &gitlab.Issue{
		IID:         iid,
		Title:       title,
		Description: description,
	}
}

// CreateMockNote creates a mock note for testing
func CreateMockNote(id int, body string) *gitlab.Note {
	return &gitlab.Note{
		ID:   id,
		Body: body,
	}
} 