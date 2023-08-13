package dev

import (
	sql "movies/sql"

	cobra "github.com/spf13/cobra"
)

func Reset() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "reset-db",
		Short: "Reset database",
	}

	// Runner
	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// Reset database
		if err := sql.Reset(ctx); err != nil {
			panic(err)
		}
	}

	return cmd
}
