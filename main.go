package main

import (
	"fmt"
	"os"

	"gitlab-manager/cmd/issues"
	"gitlab-manager/cmd/mergerequests"
	"gitlab-manager/cmd/milestones"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitlab-manager",
	Short: "A CLI tool to manage GitLab resources",
}

func init() {
	rootCmd.AddCommand(issues.IssuesCmd)
	rootCmd.AddCommand(milestones.MilestonesCmd)
	rootCmd.AddCommand(mergerequests.MergeRequestsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
