package mergerequests

import (
	"testing"

	"mpg-gitlab/cmd/utils"

	"github.com/xanzy/go-gitlab"
)

func TestAddCurrentMilestone(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient.MergeRequests

	tests := []struct {
		name        string
		projectID   int
		mrIID       int
		setupMocks  func()
		wantErr     bool
		errContains string
	}{
		{
			name:      "successfully adds milestone",
			projectID: 1,
			mrIID:    123,
			setupMocks: func() {
				// Mock MR
				mr := utils.CreateMockMR(123, "Test MR", "Fixes #456")
				mockClient.MergeRequests.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}

				// Mock milestone list
				mockClient.Milestones.ListMilestonesFunc = func(pid interface{}, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, *gitlab.Response, error) {
					return []*gitlab.Milestone{
						{ID: 1, Title: "Current", State: "active"},
					}, nil, nil
				}

				// Mock MR update
				mockClient.MergeRequests.UpdateMergeRequestFunc = func(pid interface{}, mriid int, opt *gitlab.UpdateMergeRequestOptions) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}

				// Mock issue update
				mockClient.Issues.UpdateIssueFunc = func(pid interface{}, iid int, opt *gitlab.UpdateIssueOptions) (*gitlab.Issue, *gitlab.Response, error) {
					return &gitlab.Issue{IID: iid}, nil, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "no current milestone",
			projectID: 1,
			mrIID:    123,
			setupMocks: func() {
				// Mock MR
				mr := utils.CreateMockMR(123, "Test MR", "Description")
				mockClient.MergeRequests.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}

				// Mock empty milestone list
				mockClient.Milestones.ListMilestonesFunc = func(pid interface{}, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, *gitlab.Response, error) {
					return []*gitlab.Milestone{}, nil, nil
				}
			},
			wantErr:     true,
			errContains: "no active milestone named 'Current' found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := AddCurrentMilestone(tt.projectID, tt.mrIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCurrentMilestone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("AddCurrentMilestone() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
} 