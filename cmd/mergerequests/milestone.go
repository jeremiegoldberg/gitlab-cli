package mergerequests

import (
	"fmt"
	"mpg-gitlab/cmd/utils"
	"strings"

	"github.com/xanzy/go-gitlab"
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

// AddCurrentMilestone adds the "Current" milestone to an MR and its linked issues
func AddCurrentMilestone(projectID, mrIID int) error {
	// Get the MR first
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return fmt.Errorf("failed to get merge request: %v", err)
	}

	// Find the "Current" milestone
	milestones, _, err := client.Milestones.ListMilestones(projectID, &gitlab.ListMilestonesOptions{
		State: gitlab.String("active"),
	})
	if err != nil {
		return fmt.Errorf("failed to list milestones: %v", err)
	}

	var currentMilestone *gitlab.Milestone
	for _, m := range milestones {
		if strings.EqualFold(m.Title, "current") {
			currentMilestone = m
			break
		}
	}

	if currentMilestone == nil {
		return fmt.Errorf("no active milestone named 'Current' found")
	}

	// Update MR milestone
	_, _, err = client.MergeRequests.UpdateMergeRequest(projectID, mrIID, &gitlab.UpdateMergeRequestOptions{
		MilestoneID: gitlab.Int(currentMilestone.ID),
	})
	if err != nil {
		return fmt.Errorf("failed to update merge request milestone: %v", err)
	}

	// Get linked issues
	issueIDs := utils.GetIssueIDsFromDescription(mr.Description)
	for _, issueID := range issueIDs {
		// Update each issue's milestone
		_, _, err = client.Issues.UpdateIssue(projectID, issueID, &gitlab.UpdateIssueOptions{
			MilestoneID: gitlab.Int(currentMilestone.ID),
		})
		if err != nil {
			return fmt.Errorf("failed to update issue #%d milestone: %v", issueID, err)
		}
	}

	return nil
} 