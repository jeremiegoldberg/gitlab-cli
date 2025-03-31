package mergerequests

import (
	"testing"
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