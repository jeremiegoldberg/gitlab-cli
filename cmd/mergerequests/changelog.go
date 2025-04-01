package mergerequests

import (
	"fmt"
	"mpg-gitlab/cmd/utils"
	"regexp"
	"strings"
)

const (
	// ChangelogError is returned when no valid changelog entry is found
	ChangelogError = "No valid changelog entry found"
)

// GetChangelogEntries returns changelog entries from MR and linked issues
func GetChangelogEntries(projectID, mrIID int) (string, error) {
	// Get the MR
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request: %v", err)
	}

	// Check MR description first
	entry := findChangelogEntry(mr.Description)
	if entry != "" {
		if !strings.Contains(strings.ToLower(entry), "[no-changelog-entry]") {
			return fmt.Sprintf("MR #%d: %s", mr.IID, entry), nil
		}
	}

	// Check linked issues
	issueIDs := utils.GetIssueIDsFromDescription(mr.Description)
	for _, issueID := range issueIDs {
		issue, _, err := client.Issues.GetIssue(projectID, issueID, nil)
		if err != nil {
			continue // Skip issues we can't access
		}
		entry := findChangelogEntry(issue.Description)
		if entry != "" && !strings.Contains(strings.ToLower(entry), "[no-changelog-entry]") {
			return fmt.Sprintf("#%d: %s", issue.IID, entry), nil
		}
	}

	return ChangelogError, nil
}

// findChangelogEntry extracts changelog entry from text
func findChangelogEntry(text string) string {
	// Case insensitive regex for changelog entries
	re := regexp.MustCompile(`(?i)\[(feature|improvement|fix|infra|no-changelog-entry)\].*`)
	match := re.FindString(text)
	if match == "" {
		return ""
	}
	return strings.TrimSpace(match)
}
