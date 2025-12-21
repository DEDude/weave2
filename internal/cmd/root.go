package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "weave2",
	Short: "Weave notes to RDF/graph",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: wire subcommands
	},
}
