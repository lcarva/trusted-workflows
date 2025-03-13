package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tw",
	Short: "tw (Trusted Workflows) is used to generate SLSA Provenance for different CI Providers.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
