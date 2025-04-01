package milestones

import (
	"fmt"
	"strings"

	"mpg-gitlab/cmd/mergerequests"

	"github.com/xanzy/go-gitlab"
)

// AddChangelogFromMR adds changelog entry from a single merge request to its milestone
func AddChangelogFromMR(projectID, mrIID int) error {
	// Get the MR first
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return fmt.Errorf("failed to get merge request: %v", err)
	}

	// Check if MR has milestone
	if mr.Milestone == nil {
		return fmt.Errorf("merge request #%d has no milestone assigned", mrIID)
	}

	milestoneID := mr.Milestone.ID

	// Get milestone
	milestone, _, err := client.Milestones.GetMilestone(projectID, milestoneID, nil)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %v", err)
	}

	// Get changelog entry from MR
	entry, err := mergerequests.GetChangelogEntries(projectID, mrIID)
	if err != nil {
		return fmt.Errorf("failed to get changelog entry: %v", err)
	}
	if entry == mergerequests.ChangelogError {
		return fmt.Errorf("no changelog entry found in MR #%d", mrIID)
	}

	return addEntryToMilestone(projectID, milestone, entry)
}

// AddChangelogFromMilestone adds changelog entries from all merge requests in a milestone
func AddChangelogFromMilestone(projectID, milestoneID int) error {
	// Get milestone
	milestone, _, err := client.Milestones.GetMilestone(projectID, milestoneID, nil)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %v", err)
	}

	// Get all MRs for this milestone
	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		Milestone: gitlab.String(fmt.Sprint(milestoneID)),
		State:     gitlab.String("merged"),
	})
	if err != nil {
		return fmt.Errorf("failed to list merge requests: %v", err)
	}

	// Collect changelog entries
	var entries []string
	for _, mr := range mrs {
		entry, err := mergerequests.GetChangelogEntries(projectID, mr.IID)
		if err != nil || entry == mergerequests.ChangelogError {
			continue
		}
		entries = append(entries, entry)
	}

	// Update milestone description if we have new entries
	if len(entries) > 0 {
		description := milestone.Description
		changelogSection := "\n\n## Changelog\n"
		for _, entry := range entries {
			// Check if entry already exists
			if !strings.Contains(description, entry) {
				changelogSection += "- " + entry + "\n"
			}
		}

		// Only update if we have new entries
		if !strings.Contains(description, changelogSection) {
			description += changelogSection
			_, _, err = client.Milestones.UpdateMilestone(projectID, milestoneID, &gitlab.UpdateMilestoneOptions{
				Description: gitlab.String(description),
			})
			if err != nil {
				return fmt.Errorf("failed to update milestone: %v", err)
			}
		}
	}

	return nil
}

// Helper function to add an entry to milestone description
func addEntryToMilestone(projectID int, milestone *gitlab.Milestone, entry string) error {
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
		_, _, err := client.Milestones.UpdateMilestone(projectID, milestone.ID, &gitlab.UpdateMilestoneOptions{
			Description: gitlab.String(description),
		})
		if err != nil {
			return fmt.Errorf("failed to update milestone: %v", err)
		}
	}

	return nil
}
