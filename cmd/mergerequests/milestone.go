package mergerequests

import (
	"fmt"
	"mpg-gitlab/cmd/utils"
)

// CheckMilestone verifies if the MR and its linked issues have a milestone
func CheckMilestone(projectID, mrIID int) error {
	// Get the MR
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return fmt.Errorf("failed to get merge request: %v", err)
	}

	// Check if MR has milestone
	if mr.Milestone == nil {
		return fmt.Errorf("merge request #%d has no milestone assigned", mrIID)
	}

	// Get linked issues
	issueIDs := utils.GetIssueIDsFromDescription(mr.Description)
	for _, issueID := range issueIDs {
		issue, _, err := client.Issues.GetIssue(projectID, issueID, nil)
		if err != nil {
			continue // Skip issues we can't access
		}
		if issue.Milestone == nil {
			return fmt.Errorf("linked issue #%d has no milestone assigned", issueID)
		}
		// Optionally: Check if issues have same milestone as MR
		if issue.Milestone.ID != mr.Milestone.ID {
			return fmt.Errorf("linked issue #%d has different milestone (%s) than MR (%s)", 
				issueID, issue.Milestone.Title, mr.Milestone.Title)
		}
	}

	return nil
} 