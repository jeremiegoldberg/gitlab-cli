package issues

import (
	"fmt"
	"log"
	"os"

	"gitlab-manager/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	client *gitlab.Client
	// Command groups
	IssuesCmd = &cobra.Command{
		Use:   "issues",
		Short: "Manage GitLab issues",
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Run:   runList,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get issue details",
		Run:   runGet,
	}

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		Run:   runCreate,
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing issue",
		Run:   runUpdate,
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an issue",
		Run:   runDelete,
	}

	getDescriptionCmd = &cobra.Command{
		Use:   "get-description",
		Short: "Get issue description",
		Run:   runGetDescription,
	}
)

func init() {
	client = utils.GetClient()

	// Add subcommands
	IssuesCmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd, getDescriptionCmd)

	// List flags
	listCmd.Flags().IntP("project", "p", 0, "Project ID")
	listCmd.Flags().StringP("state", "s", "", "Issue state (opened/closed)")

	// Get flags
	getCmd.Flags().IntP("project", "p", 0, "Project ID")
	getCmd.Flags().IntP("issue", "i", 0, "Issue IID")
	getCmd.MarkFlagRequired("issue")

	// Create flags
	createCmd.Flags().IntP("project", "p", 0, "Project ID")
	createCmd.Flags().StringP("title", "t", "", "Issue title")
	createCmd.Flags().StringP("description", "d", "", "Issue description")
	createCmd.Flags().StringP("labels", "l", "", "Comma-separated list of labels")
	createCmd.MarkFlagRequired("title")

	// Update flags
	updateCmd.Flags().IntP("project", "p", 0, "Project ID")
	updateCmd.Flags().IntP("issue", "i", 0, "Issue IID")
	updateCmd.Flags().StringP("title", "t", "", "New issue title")
	updateCmd.Flags().StringP("description", "d", "", "New issue description")
	updateCmd.Flags().StringP("state", "s", "", "Issue state (close/reopen)")
	updateCmd.MarkFlagRequired("issue")

	// Delete flags
	deleteCmd.Flags().IntP("project", "p", 0, "Project ID")
	deleteCmd.Flags().IntP("issue", "i", 0, "Issue IID")
	deleteCmd.MarkFlagRequired("issue")

	// Get description flags
	getDescriptionCmd.Flags().IntP("issue", "i", 0, "Issue IID")
	getDescriptionCmd.MarkFlagRequired("issue")
}

func runList(cmd *cobra.Command, args []string) {
	opts := &gitlab.ListIssuesOptions{}

	if state, _ := cmd.Flags().GetString("state"); state != "" {
		opts.State = gitlab.String(state)
	}

	// If running in CI, scope to current project
	if projectID, _ := utils.GetProjectID(cmd); projectID != 0 {
		opts.ProjectID = gitlab.Int(projectID)
	}

	issues, _, err := client.Issues.ListIssues(opts)
	if err != nil {
		log.Fatalf("Failed to list issues: %v", err)
	}

	for _, issue := range issues {
		fmt.Printf("#%d: [%s] %s\n", issue.IID, issue.State, issue.Title)
	}
}

func runGet(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	issueIID, _ := cmd.Flags().GetInt("issue")

	issue, _, err := client.Issues.GetIssue(projectID, issueIID)
	if err != nil {
		log.Fatalf("Failed to get issue: %v", err)
	}

	fmt.Printf("Issue #%d\n", issue.IID)
	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("State: %s\n", issue.State)
	fmt.Printf("Description:\n%s\n", issue.Description)
}

func runCreate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	if projectID == 0 {
		log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
	}

	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	labels, _ := cmd.Flags().GetString("labels")

	opts := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
	}

	if labels != "" {
		opts.Labels = &[]string{labels}
	}

	// Add CI metadata if available
	if os.Getenv("CI") != "" {
		ciMetadata := utils.GetCIMetadata()
		opts.Description = gitlab.String(*opts.Description + ciMetadata)
	}

	issue, _, err := client.Issues.CreateIssue(projectID, opts)
	if err != nil {
		log.Fatalf("Failed to create issue: %v", err)
	}

	fmt.Printf("Created issue #%d: %s\n", issue.IID, issue.Title)
}

func runUpdate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	issueIID, _ := cmd.Flags().GetInt("issue")

	opts := &gitlab.UpdateIssueOptions{}

	if title, _ := cmd.Flags().GetString("title"); title != "" {
		opts.Title = gitlab.String(title)
	}
	if description, _ := cmd.Flags().GetString("description"); description != "" {
		opts.Description = gitlab.String(description)
	}
	if state, _ := cmd.Flags().GetString("state"); state != "" {
		opts.StateEvent = gitlab.String(state)
	}

	issue, _, err := client.Issues.UpdateIssue(projectID, issueIID, opts)
	if err != nil {
		log.Fatalf("Failed to update issue: %v", err)
	}

	fmt.Printf("Updated issue #%d\n", issue.IID)
}

func runDelete(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	issueIID, _ := cmd.Flags().GetInt("issue")

	_, err := client.Issues.DeleteIssue(projectID, issueIID)
	if err != nil {
		log.Fatalf("Failed to delete issue: %v", err)
	}

	fmt.Printf("Deleted issue #%d\n", issueIID)
}

func runGetDescription(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	issueIID, _ := cmd.Flags().GetInt("issue")

	description, err := GetIssueDescription(projectID, issueIID)
	if err != nil {
		log.Fatalf("Failed to get issue description: %v", err)
	}

	fmt.Println(description)
}
