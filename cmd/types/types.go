package types

import (
	"time"

	"github.com/your-project/utils"
)

// Issue represents a GitLab issue
type Issue struct {
	IID         int        `json:"iid"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	Labels      []string   `json:"labels"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	WebURL      string     `json:"web_url"`
}

// MergeRequest represents a GitLab merge request
type MergeRequest struct {
	IID          int        `json:"iid"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	State        string     `json:"state"`
	SourceBranch string     `json:"source_branch"`
	TargetBranch string     `json:"target_branch"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	MergedAt     *time.Time `json:"merged_at,omitempty"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	WebURL       string     `json:"web_url"`
	MergeStatus  string     `json:"merge_status"`
	HasConflicts bool       `json:"has_conflicts"`
}

// Milestone represents a GitLab milestone
type Milestone struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	WebURL      string     `json:"web_url"`
}

// GetLinkedIssueIIDs returns the IIDs of issues referenced in the MR description
func (mr *MergeRequest) GetLinkedIssueIIDs() []int {
	return utils.GetIssueIDsFromDescription(mr.Description)
}

// GetDescription returns a formatted string representation of the MR description
func (mr *MergeRequest) GetDescription() string {
	if mr.Description == "" {
		return "No description provided"
	}
	return mr.Description
}

// GetDescription returns a formatted string representation of the issue description
func (i *Issue) GetDescription() string {
	if i.Description == "" {
		return "No description provided"
	}
	return i.Description
}
