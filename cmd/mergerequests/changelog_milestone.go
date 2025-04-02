package mergerequests

import (
	"fmt"
	"sort"
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

	// Update milestone description with sorted entries
	return addSortedEntryToMilestone(projectID, milestone, entry)
}

// addSortedEntryToMilestone adds a changelog entry to the milestone description
// organizing entries by category and ensuring uniqueness by MR ID
func addSortedEntryToMilestone(projectID int, milestone *gitlab.Milestone, entry string) error {
	description := milestone.Description
	if description == "" {
		description = "## Changelog\n"
	}

	// Extract MR ID from entry (format: "MR #123: [Category] Description" or "#123: [Category] Description")
	mrID := ""
	if strings.HasPrefix(entry, "MR #") {
		parts := strings.SplitN(entry, ":", 2)
		mrID = strings.TrimPrefix(parts[0], "MR ")
		entry = strings.TrimSpace(parts[1])
	} else if strings.HasPrefix(entry, "#") {
		parts := strings.SplitN(entry, ":", 2)
		mrID = parts[0]
		entry = strings.TrimSpace(parts[1])
	}
	if mrID == "" {
		return fmt.Errorf("invalid changelog entry format, missing MR ID: %s", entry)
	}

	// Split description into sections
	sections := map[string][]string{
		"Feature":     {},
		"Improvement": {},
		"Fix":        {},
		"Infra":      {},
	}

	// Extract existing entries by category
	lines := strings.Split(description, "\n")
	currentCategory := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "### [") {
			category := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "### [")
			currentCategory = category
		} else if strings.HasPrefix(line, "- ") && currentCategory != "" {
			// Skip if this line is from the same MR
			if strings.Contains(line, mrID) {
				continue
			}
			sections[currentCategory] = append(sections[currentCategory], line)
		}
	}

	// Add new entry to appropriate category
	for category := range sections {
		if strings.Contains(entry, "["+category+"]") {
			sections[category] = append(sections[category], fmt.Sprintf("- %s (%s)", entry, mrID))
		}
	}

	// Build new description
	var newDesc strings.Builder
	newDesc.WriteString("## Changelog\n\n")

	// Add each category and its entries
	for _, category := range []string{"Feature", "Improvement", "Fix", "Infra"} {
		entries := sections[category]
		if len(entries) > 0 {
			newDesc.WriteString(fmt.Sprintf("### [%s]\n", category))
			sort.Strings(entries) // Sort entries within category
			for _, e := range entries {
				newDesc.WriteString(e + "\n")
			}
			newDesc.WriteString("\n")
		}
	}

	// Update milestone
	_, _, err := client.Milestones.UpdateMilestone(projectID, milestone.ID, &gitlab.UpdateMilestoneOptions{
		Description: gitlab.String(newDesc.String()),
	})
	if err != nil {
		return fmt.Errorf("failed to update milestone: %v", err)
	}

	return nil
} 