package mergerequests

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"mpg-gitlab/cmd/types"
	"mpg-gitlab/cmd/issues"

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

// GetMRFromCommitMessage extracts merge request IID from a commit message
func GetMRFromCommitMessage(message string) (int, error) {
	re := regexp.MustCompile(`See merge request !(\d+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no merge request reference found in commit message")
	}
	
	mrIID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid merge request IID: %v", err)
	}
	
	return mrIID, nil
}

// GetMRFromCommit gets merge request IID from a commit ID
func GetMRFromCommit(projectID int, commitID string) (int, error) {
	// Get the commit details
	commit, _, err := client.Commits.GetCommit(projectID, commitID)
	if err != nil {
		return 0, fmt.Errorf("failed to get commit: %v", err)
	}

	// Extract MR IID from commit message
	return GetMRFromCommitMessage(commit.Message)
}
