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

// cleanDescription removes common formatting and noise from text
func cleanDescription(text string) string {
	// Remove code blocks
	codeBlock := regexp.MustCompile("```[^`]*```")
	text = codeBlock.ReplaceAllString(text, "")

	// Remove inline code
	inlineCode := regexp.MustCompile("`[^`]*`")
	text = inlineCode.ReplaceAllString(text, "")

	// Remove URLs
	urlPattern := regexp.MustCompile(`https?://\S+`)
	text = urlPattern.ReplaceAllString(text, "")

	// Replace escaped newlines with actual newlines
	text = strings.ReplaceAll(text, "\\n", "\n")

	// Remove other backslashes
	text = strings.ReplaceAll(text, "\\", "")

	// Remove extra whitespace while preserving newlines
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")

	// Remove multiple consecutive newlines
	multipleNewlines := regexp.MustCompile(`\n\s*\n+`)
	text = multipleNewlines.ReplaceAllString(text, "\n")

	// Final trim
	text = strings.TrimSpace(text)

	return text
}

// GetChangelogEntries returns changelog entries from MR and linked issues
func GetChangelogEntries(projectID, mrIID int) (string, error) {
	// Get the MR
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request: %v", err)
	}

	// Check MR description first
	entry := findChangelogEntry(cleanDescription(mr.Description))
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
		entry := findChangelogEntry(cleanDescription(issue.Description))
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
