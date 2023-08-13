package goose

import (
	sql "movies/sql"

	cobra "github.com/spf13/cobra"
)

func reset() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Roll back all migrations",
		Args:  cobra.NoArgs,
	}

	// Runner
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return sql.GooseReset()
	}

	return cmd
}