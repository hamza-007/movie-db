package goose

import (
	"strings"

	cobra "github.com/spf13/cobra"
)

func Command() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "goose",
		Short: "Database migration tool",
	}

	// Description
	cmd.Long = `
	Goose is a database migration tool.
	Manage database schema by creating incremental SQL changes.
	https://pressly.github.io/goose
	`
	cmd.Long = strings.ReplaceAll(cmd.Long, "\t", "")
	cmd.Long = strings.Replace(cmd.Long, "\n", "", 1)

	// SubCommand
	cmd.AddCommand(create())
	cmd.AddCommand(up())
	cmd.AddCommand(reset())

	return cmd
}
