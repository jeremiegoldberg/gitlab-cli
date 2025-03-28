package mergerequests

import (
	"fmt"
	"log"

	"gitlab-manager/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	client *gitlab.Client
	// Command groups
	MergeRequestsCmd = &cobra.Command{
		Use:     "mr",
		Aliases: []string{"merge-requests"},
		Short:   "Manage GitLab merge requests",
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List merge requests",
		Run:   runList,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get merge request details",
		Run:   runGet,
	}

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new merge request",
		Run:   runCreate,
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing merge request",
		Run:   runUpdate,
	}

	mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "Merge a merge request",
		Run:   runMerge,
	}

	closeCmd = &cobra.Command{
		Use:   "close",
		Short: "Close a merge request",
		Run:   runClose,
	}

	getDescriptionCmd = &cobra.Command{
		Use:   "get-description",
		Short: "Get merge request description",
		Run:   runGetDescription,
	}
)

func init() {
	client = utils.GetClient()

	// Add subcommands
	MergeRequestsCmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, mergeCmd, closeCmd, getDescriptionCmd)

	// List flags
	listCmd.Flags().IntP("project", "p", 0, "Project ID")
	listCmd.Flags().StringP("state", "s", "", "MR state (opened/closed/merged/all)")
	listCmd.Flags().StringP("target", "t", "", "Target branch")

	// Get flags
	getCmd.Flags().IntP("project", "p", 0, "Project ID")
	getCmd.Flags().IntP("mr", "m", 0, "Merge Request IID")
	getCmd.MarkFlagRequired("mr")

	// Create flags
	createCmd.Flags().IntP("project", "p", 0, "Project ID")
	createCmd.Flags().StringP("source", "s", "", "Source branch")
	createCmd.Flags().StringP("target", "t", "main", "Target branch")
	createCmd.Flags().StringP("title", "T", "", "Merge request title")
	createCmd.Flags().StringP("description", "d", "", "Merge request description")
	createCmd.Flags().BoolP("remove-source", "r", false, "Remove source branch when merged")
	createCmd.MarkFlagRequired("source")
	createCmd.MarkFlagRequired("title")

	// Update flags
	updateCmd.Flags().IntP("project", "p", 0, "Project ID")
	updateCmd.Flags().IntP("mr", "m", 0, "Merge Request IID")
	updateCmd.Flags().StringP("title", "t", "", "New title")
	updateCmd.Flags().StringP("description", "d", "", "New description")
	updateCmd.Flags().StringP("target", "T", "", "New target branch")
	updateCmd.MarkFlagRequired("mr")

	// Merge flags
	mergeCmd.Flags().IntP("project", "p", 0, "Project ID")
	mergeCmd.Flags().IntP("mr", "m", 0, "Merge Request IID")
	mergeCmd.Flags().StringP("message", "M", "", "Merge commit message")
	mergeCmd.MarkFlagRequired("mr")

	// Close flags
	closeCmd.Flags().IntP("project", "p", 0, "Project ID")
	closeCmd.Flags().IntP("mr", "m", 0, "Merge Request IID")
	closeCmd.MarkFlagRequired("mr")

	// Get description flags
	getDescriptionCmd.Flags().IntP("mr", "m", 0, "Merge Request IID")
	getDescriptionCmd.MarkFlagRequired("mr")
}

func runList(cmd *cobra.Command, args []string) {
	opts := &gitlab.ListMergeRequestsOptions{}

	if state, _ := cmd.Flags().GetString("state"); state != "" {
		opts.State = gitlab.String(state)
	}
	if target, _ := cmd.Flags().GetString("target"); target != "" {
		opts.TargetBranch = gitlab.String(target)
	}

	// If running in CI, scope to current project
	if projectID, _ := utils.GetProjectID(cmd); projectID != 0 {
		opts.ProjectID = gitlab.Int(projectID)
	}

	mrs, _, err := client.MergeRequests.ListMergeRequests(opts)
	if err != nil {
		log.Fatalf("Failed to list merge requests: %v", err)
	}

	for _, mr := range mrs {
		fmt.Printf("#%d: [%s] %s\n", mr.IID, mr.State, mr.Title)
		fmt.Printf("  %s -> %s\n", mr.SourceBranch, mr.TargetBranch)
	}
}

func runGet(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	mrIID, _ := cmd.Flags().GetInt("mr")

	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		log.Fatalf("Failed to get merge request: %v", err)
	}

	fmt.Printf("Merge Request #%d\n", mr.IID)
	fmt.Printf("Title: %s\n", mr.Title)
	fmt.Printf("State: %s\n", mr.State)
	fmt.Printf("Source: %s\n", mr.SourceBranch)
	fmt.Printf("Target: %s\n", mr.TargetBranch)
	if mr.Description != "" {
		fmt.Printf("Description:\n%s\n", mr.Description)
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	if projectID == 0 {
		log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
	}

	sourceBranch, _ := cmd.Flags().GetString("source")
	targetBranch, _ := cmd.Flags().GetString("target")
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	removeSource, _ := cmd.Flags().GetBool("remove-source")

	opts := &gitlab.CreateMergeRequestOptions{
		Title:              gitlab.String(title),
		SourceBranch:       gitlab.String(sourceBranch),
		TargetBranch:       gitlab.String(targetBranch),
		RemoveSourceBranch: gitlab.Bool(removeSource),
	}

	if description != "" {
		opts.Description = gitlab.String(description)
	}

	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, opts)
	if err != nil {
		log.Fatalf("Failed to create merge request: %v", err)
	}

	fmt.Printf("Created merge request #%d: %s\n", mr.IID, mr.Title)
}

func runUpdate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	mrIID, _ := cmd.Flags().GetInt("mr")

	opts := &gitlab.UpdateMergeRequestOptions{}

	if title, _ := cmd.Flags().GetString("title"); title != "" {
		opts.Title = gitlab.String(title)
	}
	if description, _ := cmd.Flags().GetString("description"); description != "" {
		opts.Description = gitlab.String(description)
	}
	if target, _ := cmd.Flags().GetString("target"); target != "" {
		opts.TargetBranch = gitlab.String(target)
	}

	mr, _, err := client.MergeRequests.UpdateMergeRequest(projectID, mrIID, opts)
	if err != nil {
		log.Fatalf("Failed to update merge request: %v", err)
	}

	fmt.Printf("Updated merge request #%d\n", mr.IID)
}

func runMerge(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	mrIID, _ := cmd.Flags().GetInt("mr")

	opts := &gitlab.AcceptMergeRequestOptions{}
	if message, _ := cmd.Flags().GetString("message"); message != "" {
		opts.MergeCommitMessage = gitlab.String(message)
	}

	mr, _, err := client.MergeRequests.AcceptMergeRequest(projectID, mrIID, opts)
	if err != nil {
		log.Fatalf("Failed to merge request: %v", err)
	}

	fmt.Printf("Merged request #%d\n", mr.IID)
}

func runClose(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	mrIID, _ := cmd.Flags().GetInt("mr")

	opts := &gitlab.UpdateMergeRequestOptions{
		StateEvent: gitlab.String("close"),
	}

	mr, _, err := client.MergeRequests.UpdateMergeRequest(projectID, mrIID, opts)
	if err != nil {
		log.Fatalf("Failed to close merge request: %v", err)
	}

	fmt.Printf("Closed merge request #%d\n", mr.IID)
}

func runGetDescription(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	mrIID, _ := cmd.Flags().GetInt("mr")

	description, err := GetMRDescription(projectID, mrIID)
	if err != nil {
		log.Fatalf("Failed to get merge request description: %v", err)
	}

	fmt.Println(description)
}
