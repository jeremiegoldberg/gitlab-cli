package mergerequests

import (
	"testing"

	"github.com/xanzy/go-gitlab"
	"github.com/your-project/utils"
)

func TestFindChangelogEntry(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "feature entry",
			text: "Some text\n[Feature] New awesome feature\nMore text",
			want: "[Feature] New awesome feature",
		},
		{
			name: "improvement entry",
			text: "[Improvement] Better performance",
			want: "[Improvement] Better performance",
		},
		{
			name: "fix entry",
			text: "Description\n[Fix] Fixed critical bug\n",
			want: "[Fix] Fixed critical bug",
		},
		{
			name: "infra entry",
			text: "[Infra] Updated CI pipeline",
			want: "[Infra] Updated CI pipeline",
		},
		{
			name: "no changelog needed",
			text: "[No-Changelog-Entry] Internal refactor",
			want: "[No-Changelog-Entry] Internal refactor",
		},
		{
			name: "no entry",
			text: "Just a regular description",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findChangelogEntry(tt.text); got != tt.want {
				t.Errorf("findChangelogEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMRFromCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    int
		wantErr bool
	}{
		{
			name:    "valid message",
			message: "Some commit\n\nSee merge request !16895",
			want:    16895,
			wantErr: false,
		},
		{
			name:    "message with multiple lines",
			message: "Fix bug\n\nDetailed description\nSee merge request !123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "no merge request reference",
			message: "Just a commit message",
			wantErr: true,
		},
		{
			name:    "invalid merge request ID",
			message: "See merge request !abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMRFromCommitMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMRFromCommitMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetMRFromCommitMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMRFromCommit(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient.Commits

	tests := []struct {
		name      string
		projectID int
		commitID  string
		setupMock func()
		want      int
		wantErr   bool
	}{
		{
			name:      "valid commit",
			projectID: 1,
			commitID:  "abc123",
			setupMock: func() {
				mockClient.Commits.GetCommitFunc = func(pid interface{}, sha string, opts ...gitlab.RequestOptionFunc) (*gitlab.Commit, *gitlab.Response, error) {
					return &gitlab.Commit{
						ID:      "abc123",
						Message: "Fix bug\n\nSee merge request !123",
					}, nil, nil
				}
			},
			want:    123,
			wantErr: false,
		},
		{
			name:      "commit without MR reference",
			projectID: 1,
			commitID:  "def456",
			setupMock: func() {
				mockClient.Commits.GetCommitFunc = func(pid interface{}, sha string, opts ...gitlab.RequestOptionFunc) (*gitlab.Commit, *gitlab.Response, error) {
					return &gitlab.Commit{
						ID:      "def456",
						Message: "Just a commit",
					}, nil, nil
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			got, err := GetMRFromCommit(tt.projectID, tt.commitID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMRFromCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetMRFromCommit() = %v, want %v", got, tt.want)
			}
		})
	}
} 