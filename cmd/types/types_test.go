package types

import (
	"testing"
)

func TestMergeRequest_GetDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "empty description",
			description: "",
			want:        "No description provided",
		},
		{
			name:        "with description",
			description: "Test description",
			want:        "Test description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &MergeRequest{Description: tt.description}
			if got := mr.GetDescription(); got != tt.want {
				t.Errorf("MergeRequest.GetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue_GetDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "empty description",
			description: "",
			want:        "No description provided",
		},
		{
			name:        "with description",
			description: "Test description",
			want:        "Test description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Issue{Description: tt.description}
			if got := i.GetDescription(); got != tt.want {
				t.Errorf("Issue.GetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeRequest_GetLinkedIssueIIDs(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        []int
	}{
		{
			name:        "no issues",
			description: "Just a description",
			want:        nil,
		},
		{
			name:        "single issue",
			description: "Fixes #123",
			want:        []int{123},
		},
		{
			name:        "multiple issues",
			description: "This MR fixes #123 and closes #456. See also #789",
			want:        []int{123, 456, 789},
		},
		{
			name:        "duplicate issues",
			description: "Fixes #123, closes #123",
			want:        []int{123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &MergeRequest{Description: tt.description}
			got := mr.GetLinkedIssueIIDs()
			if len(got) != len(tt.want) {
				t.Errorf("MergeRequest.GetLinkedIssueIIDs() = %v, want %v", got, tt.want)
				return
			}
			// Create a map for easier comparison
			wantMap := make(map[int]bool)
			for _, id := range tt.want {
				wantMap[id] = true
			}
			for _, id := range got {
				if !wantMap[id] {
					t.Errorf("MergeRequest.GetLinkedIssueIIDs() returned unexpected ID %v", id)
				}
			}
		})
	}
}
