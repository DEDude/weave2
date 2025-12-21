package cmd

import (
	"fmt"
	"os"

	"github.com/DeDude/weave2/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg        config.Config
	vaultFlag  string
	editorFlag string
)

const (
	envVault  = "WEAVE_VAULT"
	envEditor = "WEAVE_EDITOR"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultFlag, "vault", "", "Path to notes vault (defaults to ./notes)")
	rootCmd.PersistentFlags().StringVar(&editorFlag, "editor", "", "Editor command to record in config")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "weave2",
	Short: "Weave notes to RDF/graph",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		vault, editor := resolveInputs()
		loaded, err := config.Load(vault, editor)

		if err != nil {
			return err
		}
		cfg = loaded
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: wire subcommands
	},
}

func resolveInputs() (string, string) {
	vault := vaultFlag
	if vault == "" {
		if v, ok := os.LookupEnv(envVault); ok {
			vault = v
		}
	}
	editor := editorFlag
	if editor == "" {
		if e, ok := os.LookupEnv(envEditor); ok {
			editor = e
		}
	}
	return vault, editor
}
