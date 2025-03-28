package types

import (
	"time"

	"github.com/your-project/utils"
)

// Issue represents a GitLab issue with its core attributes
// It maps to the GitLab API issue object but includes only the fields we need
type Issue struct {
	IID         int        `json:"iid"`                 // Internal ID of the issue
	Title       string     `json:"title"`               // Issue title
	Description string     `json:"description"`         // Issue description
	State       string     `json:"state"`               // Current state (opened/closed)
	Labels      []string   `json:"labels"`              // List of labels
	CreatedAt   time.Time  `json:"created_at"`          // Creation timestamp
	UpdatedAt   time.Time  `json:"updated_at"`          // Last update timestamp
	ClosedAt    *time.Time `json:"closed_at,omitempty"` // Closing timestamp, if closed
	WebURL      string     `json:"web_url"`             // Web URL to the issue
}

// MergeRequest represents a GitLab merge request with its core attributes
// It maps to the GitLab API merge request object but includes only the fields we need
type MergeRequest struct {
	IID          int        `json:"iid"`                 // Internal ID of the MR
	Title        string     `json:"title"`               // MR title
	Description  string     `json:"description"`         // MR description
	State        string     `json:"state"`               // Current state (opened/merged/closed)
	SourceBranch string     `json:"source_branch"`       // Source branch name
	TargetBranch string     `json:"target_branch"`       // Target branch name
	CreatedAt    time.Time  `json:"created_at"`          // Creation timestamp
	UpdatedAt    time.Time  `json:"updated_at"`          // Last update timestamp
	MergedAt     *time.Time `json:"merged_at,omitempty"` // Merge timestamp, if merged
	ClosedAt     *time.Time `json:"closed_at,omitempty"` // Closing timestamp, if closed
	WebURL       string     `json:"web_url"`             // Web URL to the MR
	MergeStatus  string     `json:"merge_status"`        // Current merge status
	HasConflicts bool       `json:"has_conflicts"`       // Whether MR has conflicts
}

// Milestone represents a GitLab milestone with its core attributes
// It maps to the GitLab API milestone object but includes only the fields we need
type Milestone struct {
	ID          int        `json:"id"`                   // ID of the milestone
	Title       string     `json:"title"`                // Milestone title
	Description string     `json:"description"`          // Milestone description
	State       string     `json:"state"`                // Current state (active/closed)
	CreatedAt   time.Time  `json:"created_at"`           // Creation timestamp
	UpdatedAt   time.Time  `json:"updated_at"`           // Last update timestamp
	DueDate     *time.Time `json:"due_date,omitempty"`   // Due date, if set
	StartDate   *time.Time `json:"start_date,omitempty"` // Start date, if set
	WebURL      string     `json:"web_url"`              // Web URL to the milestone
}

// GetLinkedIssueIIDs returns the IIDs of issues referenced in the MR description
// It parses the description looking for issue references like "#123" or "fixes #456"
func (mr *MergeRequest) GetLinkedIssueIIDs() []int {
	return utils.GetIssueIDsFromDescription(mr.Description)
}

// GetDescription returns a formatted string representation of the MR description
// Returns "No description provided" if the description is empty
func (mr *MergeRequest) GetDescription() string {
	if mr.Description == "" {
		return "No description provided"
	}
	return mr.Description
}

// GetDescription returns a formatted string representation of the issue description
// Returns "No description provided" if the description is empty
func (i *Issue) GetDescription() string {
	if i.Description == "" {
		return "No description provided"
	}
	return i.Description
}
