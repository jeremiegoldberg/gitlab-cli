package issues

import (
	"encoding/json"
	"fmt"

	"mpg-gitlab/cmd/types"

	"github.com/xanzy/go-gitlab"
)

// ReadIssue gets an issue and returns it as a structured type
func ReadIssue(projectID, issueIID int) (*types.Issue, error) {
	issue, _, err := client.Issues.GetIssue(projectID, issueIID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	return ConvertGitLabIssue(issue), nil
}

// ReadIssues gets a list of issues and returns them as structured types
func ReadIssues(opts *gitlab.ListIssuesOptions) ([]types.Issue, error) {
	issues, _, err := client.Issues.ListIssues(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %v", err)
	}

	result := make([]types.Issue, len(issues))
	for i, issue := range issues {
		result[i] = *ConvertGitLabIssue(issue)
	}

	return result, nil
}

// ReadIssuesAsJSON gets issues and returns them as formatted JSON
func ReadIssuesAsJSON(opts *gitlab.ListIssuesOptions) (string, error) {
	issues, err := ReadIssues(opts)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal issues: %v", err)
	}

	return string(jsonData), nil
}

// ReadIssueAsJSON gets an issue and returns it as formatted JSON
func ReadIssueAsJSON(projectID, issueIID int) (string, error) {
	issue, err := ReadIssue(projectID, issueIID)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(issue, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal issue: %v", err)
	}

	return string(jsonData), nil
}

func ConvertGitLabIssue(issue *gitlab.Issue) *types.Issue {
	return &types.Issue{
		IID:         issue.IID,
		Title:       issue.Title,
		Description: issue.Description,
		State:       issue.State,
		Labels:      issue.Labels,
		CreatedAt:   *issue.CreatedAt,
		UpdatedAt:   *issue.UpdatedAt,
		ClosedAt:    issue.ClosedAt,
		WebURL:      issue.WebURL,
	}
}

// GetIssueDescription returns the description of an issue
func GetIssueDescription(projectID, issueIID int) (string, error) {
	issue, err := ReadIssue(projectID, issueIID)
	if err != nil {
		return "", fmt.Errorf("failed to get issue: %v", err)
	}

	return issue.GetDescription(), nil
}
