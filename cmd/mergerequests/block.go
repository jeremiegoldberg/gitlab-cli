package mergerequests

import (
	"fmt"
	"strings"
)

// IsBlocked checks if a merge request is blocked
func IsBlocked(projectID, mrIID int) (bool, error) {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get merge request: %v", err)
	}

	return strings.HasPrefix(mr.Title, "[BLOCKED]"), nil
}

// GetBlockReason returns the reason why a merge request is blocked
func GetBlockReason(projectID, mrIID int) (string, error) {
	notes, _, err := client.Notes.ListMergeRequestNotes(projectID, mrIID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get merge request notes: %v", err)
	}

	// Look for the most recent blocking note
	for i := len(notes) - 1; i >= 0; i-- {
		if strings.Contains(notes[i].Body, "🚫 **Merge Blocked**:") {
			parts := strings.SplitN(notes[i].Body, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
			return "", nil
		}
	}

	return "", nil
} 