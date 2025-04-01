package mergerequests

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

// AddChangelogToMilestone adds changelog entry from merge request to its milestone
func AddChangelogToMilestone(projectID, mrIID int) error {
	// Get the MR first
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return fmt.Errorf("failed to get merge request: %v", err)
	}

	// Check if MR has milestone
	if mr.Milestone == nil {
		return fmt.Errorf("merge request #%d has no milestone assigned", mrIID)
	}

	// Get milestone
	milestone, _, err := client.Milestones.GetMilestone(projectID, mr.Milestone.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %v", err)
	}

	// Get changelog entry from MR
	entry, err := GetChangelogEntries(projectID, mrIID)
	if err != nil {
		return fmt.Errorf("failed to get changelog entry: %v", err)
	}
	if entry == ChangelogError {
		return fmt.Errorf("no changelog entry found in MR #%d", mrIID)
	}

	// Update milestone description
	description := milestone.Description
	if description == "" {
		description = "## Changelog\n"
	}

	// Check if entry already exists
	if !strings.Contains(description, entry) {
		// Add changelog section if it doesn't exist
		if !strings.Contains(description, "## Changelog") {
			description += "\n\n## Changelog\n"
		}
		// Add the new entry
		description += "- " + entry + "\n"

		// Update milestone
		_, _, err = client.Milestones.UpdateMilestone(projectID, mr.Milestone.ID, &gitlab.UpdateMilestoneOptions{
			Description: gitlab.String(description),
		})
		if err != nil {
			return fmt.Errorf("failed to update milestone: %v", err)
		}
	}

	return nil
} 