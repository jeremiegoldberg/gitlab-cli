package utils

import (
	"regexp"
	"strconv"
)

var (
	// Matches GitLab issue references like #123, fixes #456, closes #789, etc.
	issuePattern = regexp.MustCompile(`(?i)(fixes|closes|resolves|references|refs|re|see|addresses)?\s*#(\d+)`)
)

// GetIssueIDsFromDescription extracts issue IIDs from a merge request description
func GetIssueIDsFromDescription(description string) []int {
	if description == "" {
		return nil
	}

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
func GetLinkedIssues(projectID int, description string) ([]int, error) {
	issueIDs := GetIssueIDsFromDescription(description)
	if len(issueIDs) == 0 {
		return nil, nil
	}

	return issueIDs, nil
}
