package goose

import (
	sql "movies/sql"

	cobra "github.com/spf13/cobra"
)

func up() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate the DB to the most recent version available",
		Args:  cobra.NoArgs,
	}

	// Runner
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return sql.GooseUp()
	}

	return cmd
}
