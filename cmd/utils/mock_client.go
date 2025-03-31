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

// MockMergeRequestsService provides mock implementations of GitLab MergeRequests API methods
// Each field is a function that can be customized for testing different scenarios
type MockMergeRequestsService struct {
	// GetMergeRequestFunc mocks the GetMergeRequest API call
	GetMergeRequestFunc func(pid interface{}, mergeRequest int, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	// ListMergeRequestsFunc mocks the ListMergeRequests API call
	ListMergeRequestsFunc func(opt *gitlab.ListMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	// CreateMergeRequestFunc mocks the CreateMergeRequest API call
	CreateMergeRequestFunc func(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	// UpdateMergeRequestFunc mocks the UpdateMergeRequest API call
	UpdateMergeRequestFunc func(pid interface{}, mergeRequest int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	// AcceptMergeRequestFunc mocks the AcceptMergeRequest API call
	AcceptMergeRequestFunc func(pid interface{}, mergeRequest int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
}

// MockMilestonesService provides mock implementations of GitLab Milestones API methods
// Each field is a function that can be customized for testing different scenarios
type MockMilestonesService struct {
	// GetMilestoneFunc mocks the GetMilestone API call
	GetMilestoneFunc func(pid interface{}, milestone int, options ...gitlab.RequestOptionFunc) (*gitlab.Milestone, *gitlab.Response, error)
	// ListMilestonesFunc mocks the ListMilestones API call
	ListMilestonesFunc func(pid interface{}, opt *gitlab.ListMilestonesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Milestone, *gitlab.Response, error)
	// CreateMilestoneFunc mocks the CreateMilestone API call
	CreateMilestoneFunc func(pid interface{}, opt *gitlab.CreateMilestoneOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Milestone, *gitlab.Response, error)
	// UpdateMilestoneFunc mocks the UpdateMilestone API call
	UpdateMilestoneFunc func(pid interface{}, milestone int, opt *gitlab.UpdateMilestoneOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Milestone, *gitlab.Response, error)
	// DeleteMilestoneFunc mocks the DeleteMilestone API call
	DeleteMilestoneFunc func(pid interface{}, milestone int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
}
