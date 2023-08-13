package cmd

import (
	"os"

	logger "movies/utils/logger"
	server "movies/utils/server"

	cobra "github.com/spf13/cobra"
)

func Serve() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Launch server",
	}

	// Runner
	cmd.Run = func(cmd *cobra.Command, args []string) {
		logger.Info(cmd.Context(), "Launch server at :%s", os.Getenv("PORT"))
		if err := server.Start(cmd.Context()); err != nil {
			logger.Error(cmd.Context(), err.Error())
		}
	}

	return cmd
}
