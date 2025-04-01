package milestones

import (
	"fmt"
	"strings"

	"mpg-gitlab/cmd/mergerequests"

	"github.com/xanzy/go-gitlab"
)

// AddChangelogToMilestone adds changelog entries from MRs to milestone release notes
func AddChangelogToMilestone(projectID, milestoneID int) error {
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