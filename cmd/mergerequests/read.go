package mergerequests

import (
	"encoding/json"
	"fmt"

	"gitlab-cli/cmd/types"

	"github.com/xanzy/go-gitlab"
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

	var issues []types.Issue
	for _, iid := range issueIIDs {
		issue, err := client.Issues.GetIssue(projectID, iid)
		if err != nil {
			continue // Skip issues we can't fetch
		}
		issues = append(issues, *convertGitLabIssue(issue))
	}

	return issues, nil
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
