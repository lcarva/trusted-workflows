package cmd

import (
	"fmt"

	"github.com/lcarva/trusted-workflows/internal/generator"
	"github.com/spf13/cobra"
)

var githubCommand = &cobra.Command{
	Use:   "github",
	Short: "Generate SLSA Provenance for a GitHub Workflow",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Job ID: %d\n", workflowRunId)
		fmt.Printf("Repo slug: %s%s\n", owner, repo)
		return generator.GitHubGenerate(cmd.Context(), workflowRunId, jobName, owner, repo)
	},
	SilenceUsage: true,
}

var (
	owner         string
	repo          string
	workflowRunId int
	jobName       string
)

func init() {
	githubCommand.Flags().IntVarP(&workflowRunId, "workflow-run-id", "w", 0,
		"Workflow Run ID to gather information from")
	if err := githubCommand.MarkFlagRequired("workflow-run-id"); err != nil {
		panic(err)
	}

	githubCommand.Flags().StringVarP(&jobName, "job-name", "j", "",
		"Job name withing workflow to inspect")
	// TODO: This could optional if there's only a single job in the workflow.
	if err := githubCommand.MarkFlagRequired("job-name"); err != nil {
		panic(err)
	}

	githubCommand.Flags().StringVarP(&owner, "owner", "o", "", "GitHub owner (user or organization)")
	if err := githubCommand.MarkFlagRequired("owner"); err != nil {
		panic(err)
	}

	githubCommand.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository")
	if err := githubCommand.MarkFlagRequired("repo"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(githubCommand)
}
