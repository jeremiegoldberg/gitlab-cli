package utils

import (
	"github.com/xanzy/go-gitlab"
)

// MockGitLabClient implements a mock GitLab client for testing
type MockGitLabClient struct {
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	Milestones    *MockMilestonesService
}

type MockIssuesService struct {
	GetIssueFunc    func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	ListIssuesFunc  func(opt *gitlab.ListIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
	CreateIssueFunc func(pid interface{}, opt *gitlab.CreateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	UpdateIssueFunc func(pid interface{}, issue int, opt *gitlab.UpdateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	DeleteIssueFunc func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
}

// Implement similar mock services for MergeRequests and Milestones
