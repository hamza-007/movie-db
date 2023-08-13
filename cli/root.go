package cli

import (
	cmd "movies/cli/cmd"
	dev "movies/cli/dev"
	goose "movies/cli/goose"

	lo "github.com/samber/lo"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func Execute() error {
	// Command
	root := &cobra.Command{
		Use:   "movies",
		Short: "movies backend",
		Long:  "movies backend",
	}

	// Verbose flag
	root.PersistentFlags().BoolP("verbose", "v", false, "enable verbose mode")
	lo.Must0(viper.BindPFlag("verbose", root.PersistentFlags().Lookup("verbose")))

	// SubCommand
	root.AddCommand(goose.Command())
	root.AddCommand(cmd.Serve())
	root.AddCommand(dev.Command())

	return root.Execute()
}
