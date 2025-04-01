package mergerequests

import (
	"testing"

	"mpg-gitlab/cmd/utils"

	"github.com/xanzy/go-gitlab"
)

func TestIsBlocked(t *testing.T) {
	mockClient := utils.MockClient()
	client = mockClient.MergeRequests

	tests := []struct {
		name        string
		projectID   int
		mrIID       int
		setupMocks  func()
		want        bool
		wantErr     bool
		errContains string
	}{
		{
			name:      "detects blocked MR",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				mr := utils.CreateMockMR(1, "[BLOCKED] Test MR", "Description")
				mockClient.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "detects unblocked MR",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				mr := utils.CreateMockMR(1, "Test MR", "Description")
				mockClient.GetMergeRequestFunc = func(pid interface{}, mriid int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
					return mr, nil, nil
				}
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := IsBlocked(tt.projectID, tt.mrIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBlocked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsBlocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBlockReason(t *testing.T) {
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
			name:      "finds block reason",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				notes := []*gitlab.Note{
					utils.CreateMockNote(1, "🚫 **Merge Blocked**: Needs review"),
				}
				mockClient.Notes.ListMergeRequestNotesFunc = func(pid interface{}, mriid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, *gitlab.Response, error) {
					return notes, nil, nil
				}
			},
			want:    "Needs review",
			wantErr: false,
		},
		{
			name:      "returns empty when no block reason",
			projectID: 1,
			mrIID:    1,
			setupMocks: func() {
				notes := []*gitlab.Note{
					utils.CreateMockNote(1, "Regular comment"),
				}
				mockClient.Notes.ListMergeRequestNotesFunc = func(pid interface{}, mriid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, *gitlab.Response, error) {
					return notes, nil, nil
				}
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := GetBlockReason(tt.projectID, tt.mrIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockReason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBlockReason() = %v, want %v", got, tt.want)
			}
		})
	}
} 