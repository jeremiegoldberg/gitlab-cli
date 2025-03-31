package milestones

import (
	"encoding/json"
	"fmt"
	"time"

	"gitlab-manager/cmd/types"

	"github.com/xanzy/go-gitlab"
)

// ReadMilestone gets a milestone and returns it as a structured type
func ReadMilestone(projectID, milestoneID int) (*types.Milestone, error) {
	milestone, _, err := client.Milestones.GetMilestone(projectID, milestoneID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone: %v", err)
	}

	return convertGitLabMilestone(milestone), nil
}

// ReadMilestones gets a list of milestones and returns them as structured types
func ReadMilestones(projectID int, opts *gitlab.ListMilestonesOptions) ([]types.Milestone, error) {
	milestones, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list milestones: %v", err)
	}

	result := make([]types.Milestone, len(milestones))
	for i, milestone := range milestones {
		result[i] = *convertGitLabMilestone(milestone)
	}

	return result, nil
}

// ReadMilestonesAsJSON gets milestones and returns them as formatted JSON
func ReadMilestonesAsJSON(projectID int, opts *gitlab.ListMilestonesOptions) (string, error) {
	milestones, err := ReadMilestones(projectID, opts)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(milestones, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal milestones: %v", err)
	}

	return string(jsonData), nil
}

// ReadMilestoneAsJSON gets a milestone and returns it as formatted JSON
func ReadMilestoneAsJSON(projectID, milestoneID int) (string, error) {
	milestone, err := ReadMilestone(projectID, milestoneID)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(milestone, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal milestone: %v", err)
	}

	return string(jsonData), nil
}

func isoTimeToTime(t *gitlab.ISOTime) *time.Time {
	if t == nil {
		return nil
	}
	timeVal := time.Time(*t)
	return &timeVal
}

func convertGitLabMilestone(milestone *gitlab.Milestone) *types.Milestone {
	return &types.Milestone{
		ID:          milestone.ID,
		Title:       milestone.Title,
		Description: milestone.Description,
		State:       milestone.State,
		CreatedAt:   *milestone.CreatedAt,
		UpdatedAt:   *milestone.UpdatedAt,
		DueDate:     isoTimeToTime(milestone.DueDate),
		StartDate:   isoTimeToTime(milestone.StartDate),
		WebURL:      milestone.WebURL,
	}
}
