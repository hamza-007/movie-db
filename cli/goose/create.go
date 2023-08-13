package goose

import (
	sql "movies/sql"

	cobra "github.com/spf13/cobra"
)

func create() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates new migration file with the current timestamp",
		Args:  cobra.ExactArgs(1),
	}

	// Runner
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return sql.GooseCreate(args[0])
	}

	return cmd
}
