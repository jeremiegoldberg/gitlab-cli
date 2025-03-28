package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// getProjectID tries to get project ID from flag or CI variable
func getProjectID(cmd *cobra.Command) (int, error) {
	projectID, _ := cmd.Flags().GetInt("project")
	if projectID == 0 {
		// Try to get from CI variable
		if ciProjectID := os.Getenv("CI_PROJECT_ID"); ciProjectID != "" {
			return strconv.Atoi(ciProjectID)
		}
	}
	return projectID, nil
}

var listIssuesCmd = &cobra.Command{
	Use:   "list-issues",
	Short: "List all issues",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &gitlab.ListIssuesOptions{}

		// If running in CI, scope to current project
		if ciProjectID := os.Getenv("CI_PROJECT_ID"); ciProjectID != "" {
			pid, _ := strconv.Atoi(ciProjectID)
			opts.ProjectID = gitlab.Int(pid)
		}

		issues, _, err := client.Issues.ListIssues(opts)
		if err != nil {
			log.Fatalf("Failed to list issues: %v", err)
		}

		for _, issue := range issues {
			fmt.Printf("#%d: %s\n", issue.IID, issue.Title)
		}
	},
}

var createIssueCmd = &cobra.Command{
	Use:   "create-issue",
	Short: "Create a new issue",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		projectID, _ := getProjectID(cmd)

		if projectID == 0 {
			log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
		}

		opts := &gitlab.CreateIssueOptions{
			Title:       gitlab.String(title),
			Description: gitlab.String(description),
		}

		// Add CI metadata if available
		if os.Getenv("CI") != "" {
			ciMetadata := fmt.Sprintf("\n\n---\nCreated by CI job: %s\nPipeline: %s\nBranch: %s",
				os.Getenv("CI_JOB_URL"),
				os.Getenv("CI_PIPELINE_URL"),
				os.Getenv("CI_COMMIT_REF_NAME"))
			opts.Description = gitlab.String(*opts.Description + ciMetadata)
		}

		issue, _, err := client.Issues.CreateIssue(projectID, opts)
		if err != nil {
			log.Fatalf("Failed to create issue: %v", err)
		}

		fmt.Printf("Created issue #%d: %s\n", issue.IID, issue.Title)
	},
}

var listMilestonesCmd = &cobra.Command{
	Use:   "list-milestones",
	Short: "List milestones for a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := getProjectID(cmd)
		if projectID == 0 {
			log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
		}

		milestones, _, err := client.Milestones.ListMilestones(projectID, &gitlab.ListMilestonesOptions{})
		if err != nil {
			log.Fatalf("Failed to list milestones: %v", err)
		}

		for _, milestone := range milestones {
			fmt.Printf("%s: %s\n", milestone.Title, milestone.Description)
		}
	},
}

var listMergeRequestsCmd = &cobra.Command{
	Use:   "list-mrs",
	Short: "List merge requests",
	Run: func(cmd *cobra.Command, args []string) {
		opts := &gitlab.ListMergeRequestsOptions{}

		// If running in CI, scope to current project
		if ciProjectID := os.Getenv("CI_PROJECT_ID"); ciProjectID != "" {
			pid, _ := strconv.Atoi(ciProjectID)
			opts.ProjectID = gitlab.Int(pid)
		}

		mrs, _, err := client.MergeRequests.ListMergeRequests(opts)
		if err != nil {
			log.Fatalf("Failed to list merge requests: %v", err)
		}

		for _, mr := range mrs {
			fmt.Printf("#%d: %s [%s]\n", mr.IID, mr.Title, mr.State)
		}
	},
}
