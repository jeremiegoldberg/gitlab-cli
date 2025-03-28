package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var client *gitlab.Client

func init() {
	// Initialize client
	token := os.Getenv("CI_JOB_TOKEN")
	if token == "" {
		token = os.Getenv("GITLAB_TOKEN")
	}

	var err error
	baseURL := os.Getenv("CI_API_V4_URL")
	if baseURL != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	} else {
		client, err = gitlab.NewClient(token)
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to create GitLab client: %v", err))
	}
}

func GetClient() *gitlab.Client {
	return client
}

func GetProjectID(cmd *cobra.Command) (int, error) {
	projectID, _ := cmd.Flags().GetInt("project")
	if projectID == 0 {
		if ciProjectID := os.Getenv("CI_PROJECT_ID"); ciProjectID != "" {
			return strconv.Atoi(ciProjectID)
		}
	}
	return projectID, nil
}

func GetCIMetadata() string {
	return fmt.Sprintf("\n\n---\nCreated by CI job: %s\nPipeline: %s\nBranch: %s",
		os.Getenv("CI_JOB_URL"),
		os.Getenv("CI_PIPELINE_URL"),
		os.Getenv("CI_COMMIT_REF_NAME"))
}
