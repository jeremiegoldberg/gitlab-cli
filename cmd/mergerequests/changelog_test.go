package mergerequests

import (
	"testing"
	"strings"

	"mpg-gitlab/cmd/utils"

	"github.com/xanzy/go-gitlab"
)

func TestGetChangelogEntries(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient.MergeRequests // Use the mock service directly

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

func TestAddSortedEntryToMilestone(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient.Milestones

	tests := []struct {
		name        string
		description string
		newEntry    string
		want        string
	}{
		{
			name: "add to existing category",
			description: `## Changelog

### [Feature]
- [Feature] Existing feature (#123)

### [Fix]
- [Fix] Existing fix (#456)
`,
			newEntry: "MR #789: [Feature] New feature",
			want: `## Changelog

### [Feature]
- [Feature] Existing feature (#123)
- [Feature] New feature (#789)

### [Fix]
- [Fix] Existing fix (#456)

`,
		},
		{
			name: "skip duplicate MR entry",
			description: `## Changelog

### [Feature]
- [Feature] Old feature (#123)
`,
			newEntry: "MR #123: [Feature] Updated feature",
			want: `## Changelog

### [Feature]
- [Feature] Updated feature (#123)

`,
		},
		{
			name: "create new category",
			description: `## Changelog

### [Feature]
- [Feature] Existing feature (#123)
`,
			newEntry: "#456: [Infra] New infrastructure",
			want: `## Changelog

### [Feature]
- [Feature] Existing feature (#123)

### [Infra]
- [Infra] New infrastructure (#456)

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			milestone := &gitlab.Milestone{
				ID:          1,
				Description: tt.description,
			}

			var capturedDesc string
			mockClient.Milestones.UpdateMilestoneFunc = func(pid interface{}, milestone int, opt *gitlab.UpdateMilestoneOptions) (*gitlab.Milestone, *gitlab.Response, error) {
				capturedDesc = *opt.Description
				return nil, nil, nil
			}

			err := addSortedEntryToMilestone(1, milestone, tt.newEntry)
			if err != nil {
				t.Errorf("addSortedEntryToMilestone() error = %v", err)
				return
			}

			// Compare normalized strings (remove extra whitespace)
			normalizeString := func(s string) string {
				return strings.Join(strings.Fields(s), " ")
			}
			if normalizeString(capturedDesc) != normalizeString(tt.want) {
				t.Errorf("Description mismatch\nGot:\n%s\nWant:\n%s", capturedDesc, tt.want)
			}
		})
	}
} 