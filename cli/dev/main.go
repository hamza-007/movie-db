package dev

import (
	lo "github.com/samber/lo"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func Command() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Development purpose command",
	}

	// Flags
	cmd.PersistentFlags().String("config", "config.yml", "Path to config file")

	// Required flags
	lo.Must0(viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")))

	// SubCommand
	cmd.AddCommand(Init())
	cmd.AddCommand(Reset())

	return cmd
}
