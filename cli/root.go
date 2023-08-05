package cli

import (
	goose "movies/cli/migration"

	cobra "github.com/spf13/cobra"
)

func Execute() error {
	// Command
	root := &cobra.Command{
		Use:   "hamza",
	}

	root.AddCommand(goose.Command())

	return root.Execute()
}