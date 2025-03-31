package mergerequests

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"gitlab-manager/cmd/types"
	"gitlab-manager/cmd/issues"

	"github.com/xanzy/go-gitlab"
)

const (
	changelogError = "ERROR"
	noChangelogComment = "No changelog found, please add one of the [Feature] / [Improvement] / [Fix] / [Infra] / [No-Changelog-Entry] in the Issue mentioned in the MR description"
)

// ReadMergeRequest gets a merge request and returns it as a structured type
func ReadMergeRequest(projectID, mrIID int) (*types.MergeRequest, error) {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request: %v", err)
	}

	return convertGitLabMR(mr), nil
}

// ReadMergeRequests gets a list of merge requests and returns them as structured types
func ReadMergeRequests(opts *gitlab.ListMergeRequestsOptions) ([]types.MergeRequest, error) {
	mrs, _, err := client.MergeRequests.ListMergeRequests(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list merge requests: %v", err)
	}

	result := make([]types.MergeRequest, len(mrs))
	for i, mr := range mrs {
		result[i] = *convertGitLabMR(mr)
	}

	return result, nil
}

// ReadMergeRequestsAsJSON gets merge requests and returns them as formatted JSON
func ReadMergeRequestsAsJSON(opts *gitlab.ListMergeRequestsOptions) (string, error) {
	mrs, err := ReadMergeRequests(opts)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(mrs, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal merge requests: %v", err)
	}

	return string(jsonData), nil
}

// ReadMergeRequestAsJSON gets a merge request and returns it as formatted JSON
func ReadMergeRequestAsJSON(projectID, mrIID int) (string, error) {
	mr, err := ReadMergeRequest(projectID, mrIID)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(mr, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal merge request: %v", err)
	}

	return string(jsonData), nil
}

// Helper function to convert GitLab MR to our type
func convertGitLabMR(mr *gitlab.MergeRequest) *types.MergeRequest {
	return &types.MergeRequest{
		IID:          mr.IID,
		Title:        mr.Title,
		Description:  mr.Description,
		State:        mr.State,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
		MergedAt:     mr.MergedAt,
		ClosedAt:     mr.ClosedAt,
		WebURL:       mr.WebURL,
		MergeStatus:  mr.MergeStatus,
		HasConflicts: mr.HasConflicts,
	}
}

// GetLinkedIssues returns the issues referenced in a merge request
func GetLinkedIssues(projectID, mrIID int) ([]types.Issue, error) {
	mr, err := ReadMergeRequest(projectID, mrIID)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request: %v", err)
	}

	issueIIDs := mr.GetLinkedIssueIIDs()
	if len(issueIIDs) == 0 {
		return nil, nil
	}

	var result []types.Issue
	for _, iid := range issueIIDs {
		issue, err := issues.ReadIssue(projectID, iid)
		if err != nil {
			log.Printf("Warning: Failed to get issue #%d: %v", iid, err)
			continue
		}
		result = append(result, *issue)
	}

	return result, nil
}

// GetLinkedIssuesAsJSON returns the referenced issues as JSON
func GetLinkedIssuesAsJSON(projectID, mrIID int) (string, error) {
	issues, err := GetLinkedIssues(projectID, mrIID)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal issues: %v", err)
	}

	return string(jsonData), nil
}

// GetMRDescription returns the description of a merge request
func GetMRDescription(projectID, mrIID int) (string, error) {
	mr, err := ReadMergeRequest(projectID, mrIID)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request: %v", err)
	}

	return mr.GetDescription(), nil
}

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

// Add these functions
func IsBlocked(projectID, mrIID int) (bool, error) {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get merge request: %v", err)
	}

	return strings.HasPrefix(mr.Title, "[BLOCKED]"), nil
}

func GetBlockReason(projectID, mrIID int) (string, error) {
	notes, _, err := client.Notes.ListMergeRequestNotes(projectID, mrIID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request notes: %v", err)
	}

	// Look for the most recent blocking note
	for i := len(notes) - 1; i >= 0; i-- {
		if strings.Contains(notes[i].Body, "🚫 **Merge Blocked**:") {
			parts := strings.SplitN(notes[i].Body, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
			return "", nil
		}
	}

	return "", nil
}
