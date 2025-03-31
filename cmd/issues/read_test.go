package issues

import (
	"encoding/json"
	"testing"
	"time"

	"mpg-gitlab/cmd/types"
	"mpg-gitlab/cmd/utils"

	"github.com/xanzy/go-gitlab"
)

func TestGetIssueDescription(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient // Set the mock client

	now := time.Now()
	tests := []struct {
		name      string
		projectID int
		issueIID  int
		mockIssue *gitlab.Issue
		want      string
		wantErr   bool
	}{
		{
			name:      "issue with description",
			projectID: 1,
			issueIID:  123,
			mockIssue: &gitlab.Issue{
				IID:         123,
				Description: "Test description",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			want:    "Test description",
			wantErr: false,
		},
		{
			name:      "issue without description",
			projectID: 1,
			issueIID:  123,
			mockIssue: &gitlab.Issue{
				IID:         123,
				Description: "",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			want:    "No description provided",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock response
			mockClient.Issues.GetIssueFunc = func(pid interface{}, iid int, opt *gitlab.GetIssueOptions) (*gitlab.Issue, *gitlab.Response, error) {
				return tt.mockIssue, nil, nil
			}

			got, err := GetIssueDescription(tt.projectID, tt.issueIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIssueDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIssueDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadIssuesAsJSON(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		issues  []types.Issue
		want    string
		wantErr bool
	}{
		{
			name: "single issue",
			issues: []types.Issue{
				{
					IID:         123,
					Title:       "Test Issue",
					Description: "Test Description",
					State:       "opened",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create expected JSON
			expected, _ := json.MarshalIndent(tt.issues, "", "  ")
			tt.want = string(expected)

			// TODO: Mock GitLab client
			// This requires setting up a mock client implementation
			got, err := ReadIssuesAsJSON(&gitlab.ListIssuesOptions{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadIssuesAsJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadIssuesAsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
