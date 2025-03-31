package mergerequests

import (
	"testing"

	"mpg-gitlab/cmd/utils"

	"github.com/xanzy/go-gitlab"
)

func TestGetChangelogEntries(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient

	tests := []struct {
		name        string
		projectID   int
		mrIID       int
		setupMocks  func()
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:      "finds changelog in linked issue",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				// Setup MR with linked issue
				mr := utils.CreateMockMR(1, "Test MR", "Fixes #1")
				mockClient.MergeRequests.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}

				// Setup issue with changelog
				issue := utils.CreateMockIssue(1, "Test Issue", "[Feature] New feature")
				mockClient.Issues.GetIssueFunc = func(pid interface{}, iid int, opts ...gitlab.RequestOptionFunc) (*gitlab.Issue, *gitlab.Response, error) {
					return issue, nil, nil
				}
			},
			want:    "#1: [Feature] New feature",
			wantErr: false,
		},
		{
			name:      "finds changelog in MR when no linked issues",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				// Setup MR with changelog
				mr := utils.CreateMockMR(1, "Test MR", "[Fix] Bug fix")
				mockClient.MergeRequests.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}
			},
			want:    "MR #1: [Fix] Bug fix",
			wantErr: false,
		},
		{
			name:      "returns error when no changelog found",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				// Setup MR with no changelog
				mr := utils.CreateMockMR(1, "Test MR", "No changelog here")
				mockClient.MergeRequests.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}
			},
			want:    changelogError,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := GetChangelogEntries(tt.projectID, tt.mrIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChangelogEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChangelogEntries() = %v, want %v", got, tt.want)
			}
		})
	}
} 