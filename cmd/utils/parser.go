package utils

import (
	"regexp"
	"sort"
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
	issuePattern = regexp.MustCompile(`(?i)(?:fixes|closes|resolves|references|refs|re|see|addresses)?\s*#(\d+)`)
)

// GetIssueIDsFromDescription extracts issue IDs from text
func GetIssueIDsFromDescription(text string) []int {
	re := regexp.MustCompile(`(?i)(?:fixes|closes|resolves|implements|addresses|re|see|refs?|#)\s*#?(\d+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	seen := make(map[int]bool)
	var result []int

	for _, match := range matches {
		if id, err := strconv.Atoi(match[1]); err == nil {
			if !seen[id] {
				seen[id] = true
				result = append(result, id)
			}
		}
	}

	// Sort the results for consistent order
	sort.Ints(result)
	return result
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
