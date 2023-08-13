package dev

import (
	sql "movies/sql"

	cobra "github.com/spf13/cobra"
)

func Init() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "init-db",
		Short: "Init database",
	}

	// Runner
	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if err := sql.Init(ctx); err != nil {
			panic(err)
		}
	}

	return cmd
}
