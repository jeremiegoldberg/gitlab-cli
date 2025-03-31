package main

import (
	"fmt"
	"os"

	"mpg-gitlab/cmd/issues"
	"mpg-gitlab/cmd/mergerequests"
	"mpg-gitlab/cmd/milestones"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mpg-gitlab",
		Short: "A GitLab CLI tool for managing merge requests",
		Long: `A command-line interface for GitLab that helps manage
merge requests, issues, and milestones with features like
changelog validation and merge blocking.`,
	}
)

func init() {
	// Add command groups
	rootCmd.AddCommand(
		mergerequests.MergeRequestsCmd,
		issues.IssuesCmd,
		milestones.MilestonesCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
