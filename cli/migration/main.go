package migration

import (
	// "fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

func Command() *cobra.Command {
	// Command
	cmd := &cobra.Command{
		Use: "create",
	}

	cmd.Flags().String("name", "", "filename")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flags().Lookup("name").Value.String()
		f, err := os.Create(filename+".go")
		if err != nil {
			return err
		}

		_, err = f.WriteString("package " + filename)
		return err
	}

	return cmd
}
