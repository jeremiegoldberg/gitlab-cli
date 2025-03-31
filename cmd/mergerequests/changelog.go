package mergerequests

import (
	"fmt"
	"regexp"
	"strings"
)

// GetChangelogEntries checks for changelog entries in MR and linked issues
func GetChangelogEntries(projectID, mrIID int) (string, error) {
	// First check linked issues
	issues, err := GetLinkedIssues(projectID, mrIID)
	if err != nil {
		return "", fmt.Errorf("failed to get linked issues: %v", err)
	}

	// If we have linked issues, check their descriptions
	if len(issues) > 0 {
		for _, issue := range issues {
			entry := findChangelogEntry(issue.Description)
			if entry != "" {
				return fmt.Sprintf("#%d: %s", issue.IID, entry), nil
			}
		}
		// If no changelog found in any linked issue
		return changelogError, nil
	}

	// If no linked issues, check MR description
	mr, err := ReadMergeRequest(projectID, mrIID)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request: %v", err)
	}

	entry := findChangelogEntry(mr.Description)
	if entry != "" {
		return fmt.Sprintf("MR #%d: %s", mr.IID, entry), nil
	}

	return changelogError, nil
}

// findChangelogEntry searches for changelog entries in text
func findChangelogEntry(text string) string {
	// Define changelog patterns
	patterns := []string{
		`\[Feature\].*`,
		`\[Improvement\].*`,
		`\[Fix\].*`,
		`\[Infra\].*`,
		`\[No-Changelog-Entry\].*`,
	}

	// Check each line for changelog entries
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		for _, pattern := range patterns {
			if matched, err := regexp.MatchString(pattern, line); err == nil && matched {
				return strings.TrimSpace(line)
			}
		}
	}
	return ""
} 