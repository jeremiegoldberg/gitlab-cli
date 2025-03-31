package utils

import (
	"github.com/xanzy/go-gitlab"
)

// MockGitLabClient implements a mock GitLab client for testing
type MockGitLabClient struct {
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	Milestones    *MockMilestonesService
	Notes         *MockNotesService
}

// MockIssuesService implements mock GitLab Issues API methods
type MockIssuesService struct {
	GetIssueFunc        func(pid interface{}, iid int, opts ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error)
	ListIssuesFunc      func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, *gitlab.Response, error)
	CreateIssueFunc     func(pid interface{}, opt *gitlab.CreateIssueOptions) (*gitlab.Issue, *gitlab.Response, error)
	UpdateIssueFunc     func(pid interface{}, iid int, opt *gitlab.UpdateIssueOptions) (*gitlab.Issue, *gitlab.Response, error)
	DeleteIssueFunc     func(pid interface{}, iid int) (*gitlab.Response, error)
}

// MockMergeRequestsService implements mock GitLab MergeRequests API methods
type MockMergeRequestsService struct {
	GetMergeRequestFunc           func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	ListMergeRequestsFunc         func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	CreateMergeRequestFunc        func(pid interface{}, opt *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
	UpdateMergeRequestFunc        func(pid interface{}, mriid int, opt *gitlab.UpdateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
	AcceptMergeRequestFunc        func(pid interface{}, mriid int, opt *gitlab.AcceptMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error)
	ListProjectMergeRequestsFunc  func(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error)
}

// MockMilestonesService implements mock GitLab Milestones API methods
type MockMilestonesService struct {
	GetMilestoneFunc        func(pid interface{}, milestone int, opts ...gitlab.RequestOptionFunc) (*gitlab.Milestone, *gitlab.Response, error)
	ListMilestonesFunc      func(pid interface{}, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, *gitlab.Response, error)
	CreateMilestoneFunc     func(pid interface{}, opt *gitlab.CreateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error)
	UpdateMilestoneFunc     func(pid interface{}, milestone int, opt *gitlab.UpdateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error)
	DeleteMilestoneFunc     func(pid interface{}, milestone int) (*gitlab.Response, error)
}

// MockNotesService implements mock GitLab Notes API methods
type MockNotesService struct {
	CreateMergeRequestNoteFunc    func(pid interface{}, mriid int, opt *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, *gitlab.Response, error)
	ListMergeRequestNotesFunc     func(pid interface{}, mriid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, *gitlab.Response, error)
}

// Implement mock methods for IssuesService
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

// Implement mock methods for MergeRequestsService
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
