package utils

import (
	"github.com/xanzy/go-gitlab"
)

// MockGitLabClient implements a mock GitLab client for testing
// It provides mock implementations of GitLab API services
type MockGitLabClient struct {
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	Milestones    *MockMilestonesService
}

// MockIssuesService provides mock implementations of GitLab Issues API methods
// Each field is a function that can be customized for testing different scenarios
type MockIssuesService struct {
	// GetIssueFunc mocks the GetIssue API call
	GetIssueFunc func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	// ListIssuesFunc mocks the ListIssues API call
	ListIssuesFunc func(opt *gitlab.ListIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
	// CreateIssueFunc mocks the CreateIssue API call
	CreateIssueFunc func(pid interface{}, opt *gitlab.CreateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	// UpdateIssueFunc mocks the UpdateIssue API call
	UpdateIssueFunc func(pid interface{}, issue int, opt *gitlab.UpdateIssueOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	// DeleteIssueFunc mocks the DeleteIssue API call
	DeleteIssueFunc func(pid interface{}, issue int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
}

// Implement similar mock services for MergeRequests and Milestones
