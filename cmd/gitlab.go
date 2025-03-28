package cmd

import (
	"fmt"

	"github.com/lcarva/trusted-workflows/internal/generator"
	"github.com/spf13/cobra"
)

var gitLabCommand = &cobra.Command{
	Use:   "gitlab",
	Short: "Generate SLSA Provenance for a GitLab Pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Pipelien ID: %d\n", glPipelineId)
		fmt.Printf("Repo slug: %s%s\n", glOwner, glRepo)
		return generator.GitLabGenerate(cmd.Context(), glPipelineId, glOwner, glRepo)
	},
	SilenceUsage: true,
}

var (
	glOwner      string
	glRepo       string
	glPipelineId int
)

func init() {
	gitLabCommand.Flags().IntVarP(&glPipelineId, "pipeline-id", "p", 0,
		"Pipeline ID to gather information from")
	if err := gitLabCommand.MarkFlagRequired("pipeline-id"); err != nil {
		panic(err)
	}

	gitLabCommand.Flags().StringVarP(&glOwner, "owner", "o", "", "GitLab owner (user or group)")
	if err := gitLabCommand.MarkFlagRequired("owner"); err != nil {
		panic(err)
	}

	gitLabCommand.Flags().StringVarP(&glRepo, "repo", "r", "", "GitLab repository")
	if err := gitLabCommand.MarkFlagRequired("repo"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(gitLabCommand)
}
