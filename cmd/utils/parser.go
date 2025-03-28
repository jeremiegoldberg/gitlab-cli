package utils

import (
	"regexp"
	"strconv"
)

var (
	// issuePattern matches GitLab issue references in text
	// Supports formats like:
	// - #123
	// - fixes #456
	// - closes #789
	// - resolves #101
	// - references #202
	// - refs #303
	// - re #404
	// - see #505
	// - addresses #606
	// Case insensitive to match variations like "Fixes" or "CLOSES"
	issuePattern = regexp.MustCompile(`(?i)(fixes|closes|resolves|references|refs|re|see|addresses)?\s*#(\d+)`)
)

// GetIssueIDsFromDescription extracts issue IIDs from a merge request description
// It returns a deduplicated list of issue IDs found in the text
func GetIssueIDsFromDescription(description string) []int {
	if description == "" {
		return nil
	}

	// Find all matches in the description
	matches := issuePattern.FindAllStringSubmatch(description, -1)
	if len(matches) == 0 {
		return nil
	}

	// Use map to deduplicate IDs
	issueMap := make(map[int]bool)
	for _, match := range matches {
		if id, err := strconv.Atoi(match[2]); err == nil {
			issueMap[id] = true
		}
	}

	// Convert map to slice
	var issues []int
	for id := range issueMap {
		issues = append(issues, id)
	}

	return issues
}

// GetLinkedIssues returns the Issue objects for issues referenced in a description
// It takes a project ID and description text, and returns the list of issue IDs found
func GetLinkedIssues(projectID int, description string) ([]int, error) {
	issueIDs := GetIssueIDsFromDescription(description)
	if len(issueIDs) == 0 {
		return nil, nil
	}

	return issueIDs, nil
}
