package milestones

import (
	"fmt"
	"log"
	"time"

	"mpg-gitlab/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	client *gitlab.Client
	// Command groups
	MilestonesCmd = &cobra.Command{
		Use:   "milestones",
		Short: "Manage GitLab milestones",
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List milestones",
		Run:   runList,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get milestone details",
		Run:   runGet,
	}

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new milestone",
		Run:   runCreate,
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing milestone",
		Run:   runUpdate,
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a milestone",
		Run:   runDelete,
	}

	addChangelogCmd = &cobra.Command{
		Use:   "add-changelog",
		Short: "Add changelog entries from merge requests to milestone release notes",
		Run:   runAddChangelog,
	}
)

func init() {
	client = utils.GetClient()

	// Add subcommands
	MilestonesCmd.AddCommand(listCmd, getCmd, createCmd, updateCmd, deleteCmd, addChangelogCmd)

	// List flags
	listCmd.Flags().IntP("project", "p", 0, "Project ID")
	listCmd.Flags().StringP("state", "s", "", "Milestone state (active/closed)")

	// Get flags
	getCmd.Flags().IntP("project", "p", 0, "Project ID")
	getCmd.Flags().IntP("milestone", "m", 0, "Milestone ID")
	getCmd.MarkFlagRequired("milestone")

	// Create flags
	createCmd.Flags().IntP("project", "p", 0, "Project ID")
	createCmd.Flags().StringP("title", "t", "", "Milestone title")
	createCmd.Flags().StringP("description", "d", "", "Milestone description")
	createCmd.Flags().StringP("due-date", "D", "", "Due date (YYYY-MM-DD)")
	createCmd.MarkFlagRequired("title")

	// Update flags
	updateCmd.Flags().IntP("project", "p", 0, "Project ID")
	updateCmd.Flags().IntP("milestone", "m", 0, "Milestone ID")
	updateCmd.Flags().StringP("title", "t", "", "New milestone title")
	updateCmd.Flags().StringP("description", "d", "", "New milestone description")
	updateCmd.Flags().StringP("due-date", "D", "", "Due date (YYYY-MM-DD)")
	updateCmd.Flags().StringP("state", "s", "", "State event (close/activate)")
	updateCmd.MarkFlagRequired("milestone")

	// Delete flags
	deleteCmd.Flags().IntP("project", "p", 0, "Project ID")
	deleteCmd.Flags().IntP("milestone", "m", 0, "Milestone ID")
	deleteCmd.MarkFlagRequired("milestone")

	// Add changelog flags
	addChangelogCmd.Flags().IntP("merge-request", "r", 0, "Merge request IID")
	addChangelogCmd.Flags().IntP("milestone", "m", 0, "Milestone ID")
	addChangelogCmd.Flags().IntP("project", "p", 0, "Project ID")
	// Make one of them required
	addChangelogCmd.MarkFlagsMutuallyExclusive("merge-request", "milestone")
}

func stringToISOTime(date string) *gitlab.ISOTime {
	if date == "" {
		return nil
	}
	// Parse the date string first
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	// Convert time.Time to ISOTime
	isoTime := gitlab.ISOTime(parsed)
	return &isoTime
}

func runList(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	if projectID == 0 {
		log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
	}

	opts := &gitlab.ListMilestonesOptions{}
	if state, _ := cmd.Flags().GetString("state"); state != "" {
		opts.State = gitlab.String(state)
	}

	milestones, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		log.Fatalf("Failed to list milestones: %v", err)
	}

	for _, milestone := range milestones {
		fmt.Printf("#%d: [%s] %s\n", milestone.ID, milestone.State, milestone.Title)
		if milestone.DueDate != nil {
			fmt.Printf("  Due: %s\n", milestone.DueDate.String())
		}
	}
}

func runGet(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	milestoneID, _ := cmd.Flags().GetInt("milestone")

	milestone, _, err := client.Milestones.GetMilestone(projectID, milestoneID)
	if err != nil {
		log.Fatalf("Failed to get milestone: %v", err)
	}

	fmt.Printf("Milestone #%d\n", milestone.ID)
	fmt.Printf("Title: %s\n", milestone.Title)
	fmt.Printf("State: %s\n", milestone.State)
	if milestone.DueDate != nil {
		fmt.Printf("Due Date: %s\n", milestone.DueDate.String())
	}
	if milestone.Description != "" {
		fmt.Printf("Description:\n%s\n", milestone.Description)
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	if projectID == 0 {
		log.Fatal("Project ID is required. Provide --project flag or run in GitLab CI")
	}

	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	dueDate, _ := cmd.Flags().GetString("due-date")

	opts := &gitlab.CreateMilestoneOptions{
		Title: gitlab.String(title),
	}

	if description != "" {
		opts.Description = gitlab.String(description)
	}
	if dueDate != "" {
		opts.DueDate = stringToISOTime(dueDate)
	}

	milestone, _, err := client.Milestones.CreateMilestone(projectID, opts)
	if err != nil {
		log.Fatalf("Failed to create milestone: %v", err)
	}

	fmt.Printf("Created milestone #%d: %s\n", milestone.ID, milestone.Title)
}

func runUpdate(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	milestoneID, _ := cmd.Flags().GetInt("milestone")

	opts := &gitlab.UpdateMilestoneOptions{}

	if title, _ := cmd.Flags().GetString("title"); title != "" {
		opts.Title = gitlab.String(title)
	}
	if description, _ := cmd.Flags().GetString("description"); description != "" {
		opts.Description = gitlab.String(description)
	}
	if dueDate, _ := cmd.Flags().GetString("due-date"); dueDate != "" {
		opts.DueDate = stringToISOTime(dueDate)
	}
	if state, _ := cmd.Flags().GetString("state"); state != "" {
		opts.StateEvent = gitlab.String(state)
	}

	milestone, _, err := client.Milestones.UpdateMilestone(projectID, milestoneID, opts)
	if err != nil {
		log.Fatalf("Failed to update milestone: %v", err)
	}

	fmt.Printf("Updated milestone #%d\n", milestone.ID)
}

func runDelete(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)
	milestoneID, _ := cmd.Flags().GetInt("milestone")

	_, err := client.Milestones.DeleteMilestone(projectID, milestoneID)
	if err != nil {
		log.Fatalf("Failed to delete milestone: %v", err)
	}

	fmt.Printf("Deleted milestone #%d\n", milestoneID)
}

func runAddChangelog(cmd *cobra.Command, args []string) {
	projectID, _ := utils.GetProjectID(cmd)

	// Check which flag was provided
	if mrIID, _ := cmd.Flags().GetInt("merge-request"); mrIID != 0 {
		if err := AddChangelogFromMR(projectID, mrIID); err != nil {
			log.Fatalf("Failed to add changelog: %v", err)
		}
	} else if milestoneID, _ := cmd.Flags().GetInt("milestone"); milestoneID != 0 {
		if err := AddChangelogFromMilestone(projectID, milestoneID); err != nil {
			log.Fatalf("Failed to add changelog: %v", err)
		}
	} else {
		log.Fatal("Either --merge-request or --milestone flag is required")
	}

	fmt.Println("Successfully updated milestone changelog")
}
