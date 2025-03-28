package utils

import (
	"reflect"
	"testing"
)

// TestGetIssueIDsFromDescription tests the issue ID extraction from text
// It verifies various formats of issue references are correctly parsed
func TestGetIssueIDsFromDescription(t *testing.T) {
	tests := []struct {
		name        string // Test case name
		description string // Input text
		want        []int  // Expected issue IDs
	}{
		{
			name:        "empty description",
			description: "",
			want:        nil,
		},
		{
			name:        "no issues",
			description: "Just a description",
			want:        nil,
		},
		{
			name:        "simple reference",
			description: "#123",
			want:        []int{123},
		},
		{
			name:        "fixes keyword",
			description: "fixes #123",
			want:        []int{123},
		},
		{
			name:        "closes keyword",
			description: "closes #123",
			want:        []int{123},
		},
		{
			name:        "multiple issues",
			description: "This fixes #123 and closes #456",
			want:        []int{123, 456},
		},
		{
			name:        "with other text",
			description: "Implementation of feature request #123\nFixes bug #456",
			want:        []int{123, 456},
		},
		{
			name:        "duplicate references",
			description: "Fixes #123, also fixes #123",
			want:        []int{123},
		},
		{
			name:        "various keywords",
			description: "resolves #123\nrefs #456\nre #789\nsee #101\naddresses #202",
			want:        []int{123, 456, 789, 101, 202},
		},
		{
			name:        "case insensitive",
			description: "FIXES #123 and Closes #456",
			want:        []int{123, 456},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIssueIDsFromDescription(tt.description)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIssueIDsFromDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetLinkedIssues tests the GetLinkedIssues function
// It verifies that issue IDs are correctly extracted and returned
func TestGetLinkedIssues(t *testing.T) {
	tests := []struct {
		name        string // Test case name
		projectID   int    // Project ID input
		description string // Description input
		want        []int  // Expected issue IDs
		wantErr     bool   // Whether an error is expected
	}{
		{
			name:        "empty description",
			projectID:   1,
			description: "",
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "with issues",
			projectID:   1,
			description: "Fixes #123 and #456",
			want:        []int{123, 456},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLinkedIssues(tt.projectID, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkedIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkedIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}
