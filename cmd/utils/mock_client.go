// Package utils provides utility functions and mock implementations for testing
package utils

import (
	"github.com/xanzy/go-gitlab"
)

// MockGitLabClient implements a mock GitLab client for testing.
// It provides mock implementations of the GitLab API services:
// - Issues service for managing GitLab issues
// - MergeRequests service for managing merge requests
// - Milestones service for managing milestones
// - Notes service for managing comments and notes
type MockGitLabClient struct {
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	Milestones    *MockMilestonesService
	Notes         *MockNotesService
}

// MockClient creates a new mock GitLab client for testing.
// Returns a MockGitLabClient with initialized mock services.
// Each service can be customized by setting its Func fields.
func MockClient() *MockGitLabClient {
	return &MockGitLabClient{
		Issues:        &MockIssuesService{},
		MergeRequests: &MockMergeRequestsService{},
		Milestones:    &MockMilestonesService{},
		Notes:         &MockNotesService{},
	}
}

// MockIssuesService implements mock GitLab Issues API methods.
// Each method can be customized by setting the corresponding Func field.
// Available methods:
// - GetIssue: Get a single issue
// - ListIssues: List all issues
// - CreateIssue: Create a new issue
// - UpdateIssue: Update an existing issue
// - DeleteIssue: Delete an issue
type MockIssuesService struct {
	GetIssueFunc        func(pid interface{}, iid int, opts ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	ListIssuesFunc      func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, *gitlab.Response, error)
	CreateIssueFunc     func(pid interface{}, opt *gitlab.CreateIssueOptions) (*gitlab.Issue, *gitlab.Response, error)
	UpdateIssueFunc     func(pid interface{}, iid int, opt *gitlab.UpdateIssueOptions) (*gitlab.Issue, *gitlab.Response, error)
	DeleteIssueFunc     func(pid interface{}, iid int) (*gitlab.Response, error)
}

// MockMergeRequestsService implements mock GitLab MergeRequests API methods.
// Each method can be customized by setting the corresponding Func field.
// Available methods:
// - GetMergeRequest: Get a single merge request
// - ListMergeRequests: List all merge requests
// - ListProjectMergeRequests: List merge requests in a project
// - CreateMergeRequest: Create a new merge request
// - UpdateMergeRequest: Update an existing merge request
// - AcceptMergeRequest: Accept/merge a merge request
type MockMergeRequestsService struct {
	GetMergeRequestFunc           func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	ListMergeRequestsFunc         func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	ListProjectMergeRequestsFunc  func(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	CreateMergeRequestFunc        func(pid interface{}, opt *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
	UpdateMergeRequestFunc        func(pid interface{}, mriid int, opt *gitlab.UpdateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
	AcceptMergeRequestFunc        func(pid interface{}, mriid int, opt *gitlab.AcceptMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
}

// MockMilestonesService implements mock GitLab Milestones API methods.
// Each method can be customized by setting the corresponding Func field.
// Available methods:
// - GetMilestone: Get a single milestone
// - ListMilestones: List all milestones
// - CreateMilestone: Create a new milestone
// - UpdateMilestone: Update an existing milestone
// - DeleteMilestone: Delete a milestone
type MockMilestonesService struct {
	GetMilestoneFunc        func(pid interface{}, milestone int, opts ...gitlab.RequestOptionFunc) (*gitlab.Milestone, *gitlab.Response, error)
	ListMilestonesFunc      func(pid interface{}, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, *gitlab.Response, error)
	CreateMilestoneFunc     func(pid interface{}, opt *gitlab.CreateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error)
	UpdateMilestoneFunc     func(pid interface{}, milestone int, opt *gitlab.UpdateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error)
	DeleteMilestoneFunc     func(pid interface{}, milestone int) (*gitlab.Response, error)
}

// MockNotesService implements mock GitLab Notes API methods.
// Each method can be customized by setting the corresponding Func field.
// Available methods:
// - CreateMergeRequestNote: Create a note on a merge request
// - ListMergeRequestNotes: List all notes on a merge request
type MockNotesService struct {
	CreateMergeRequestNoteFunc    func(pid interface{}, mriid int, opt *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, *gitlab.Response, error)
	ListMergeRequestNotesFunc     func(pid interface{}, mriid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, *gitlab.Response, error)
}

// GetIssue implements the mock method
func (m *MockIssuesService) GetIssue(pid interface{}, iid int, opts ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
	if m.GetIssueFunc != nil {
		return m.GetIssueFunc(pid, iid, opts...)
	}
	return nil, nil, nil
}

func (m *MockIssuesService) ListIssues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, *gitlab.Response, error) {
	if m.ListIssuesFunc != nil {
		return m.ListIssuesFunc(opt)
	}
	return nil, nil, nil
}

// GetMergeRequest implements the mock method
func (m *MockMergeRequestsService) GetMergeRequest(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.GetMergeRequestFunc != nil {
		return m.GetMergeRequestFunc(pid, mriid, opts...)
	}
	return nil, nil, nil
}

func (m *MockMergeRequestsService) ListMergeRequests(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.ListMergeRequestsFunc != nil {
		return m.ListMergeRequestsFunc(opt)
	}
	return nil, nil, nil
}

func (m *MockMergeRequestsService) ListProjectMergeRequests(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.ListProjectMergeRequestsFunc != nil {
		return m.ListProjectMergeRequestsFunc(pid, opt)
	}
	return nil, nil, nil
}

func (m *MockMergeRequestsService) UpdateMergeRequest(pid interface{}, mriid int, opt *gitlab.UpdateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.UpdateMergeRequestFunc != nil {
		return m.UpdateMergeRequestFunc(pid, mriid, opt)
	}
	return nil, nil, nil
}

// Implement mock methods for NotesService
func (m *MockNotesService) CreateMergeRequestNote(pid interface{}, mriid int, opt *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, *gitlab.Response, error) {
	if m.CreateMergeRequestNoteFunc != nil {
		return m.CreateMergeRequestNoteFunc(pid, mriid, opt)
	}
	return nil, nil, nil
}

func (m *MockNotesService) ListMergeRequestNotes(pid interface{}, mriid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, *gitlab.Response, error) {
	if m.ListMergeRequestNotesFunc != nil {
		return m.ListMergeRequestNotesFunc(pid, mriid, opt)
	}
	return nil, nil, nil
}

// Add these methods to MockMergeRequestsService
func (m *MockMergeRequestsService) CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.CreateMergeRequestFunc != nil {
		return m.CreateMergeRequestFunc(pid, opt)
	}
	return nil, nil, nil
}

func (m *MockMergeRequestsService) AcceptMergeRequest(pid interface{}, mriid int, opt *gitlab.AcceptMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error) {
	if m.AcceptMergeRequestFunc != nil {
		return m.AcceptMergeRequestFunc(pid, mriid, opt)
	}
	return nil, nil, nil
}

// Add these methods to MockIssuesService
func (m *MockIssuesService) CreateIssue(pid interface{}, opt *gitlab.CreateIssueOptions) (*gitlab.Issue, *gitlab.Response, error) {
	if m.CreateIssueFunc != nil {
		return m.CreateIssueFunc(pid, opt)
	}
	return nil, nil, nil
}

func (m *MockIssuesService) UpdateIssue(pid interface{}, iid int, opt *gitlab.UpdateIssueOptions) (*gitlab.Issue, *gitlab.Response, error) {
	if m.UpdateIssueFunc != nil {
		return m.UpdateIssueFunc(pid, iid, opt)
	}
	return nil, nil, nil
}

func (m *MockIssuesService) DeleteIssue(pid interface{}, iid int) (*gitlab.Response, error) {
	if m.DeleteIssueFunc != nil {
		return m.DeleteIssueFunc(pid, iid)
	}
	return nil, nil
}

// Add these methods to MockMilestonesService
func (m *MockMilestonesService) CreateMilestone(pid interface{}, opt *gitlab.CreateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error) {
	if m.CreateMilestoneFunc != nil {
		return m.CreateMilestoneFunc(pid, opt)
	}
	return nil, nil, nil
}

func (m *MockMilestonesService) UpdateMilestone(pid interface{}, milestone int, opt *gitlab.UpdateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error) {
	if m.UpdateMilestoneFunc != nil {
		return m.UpdateMilestoneFunc(pid, milestone, opt)
	}
	return nil, nil, nil
}

func (m *MockMilestonesService) DeleteMilestone(pid interface{}, milestone int) (*gitlab.Response, error) {
	if m.DeleteMilestoneFunc != nil {
		return m.DeleteMilestoneFunc(pid, milestone)
	}
	return nil, nil
}
