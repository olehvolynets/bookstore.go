package cli

import (
	"github.com/spf13/cobra"

	"bookstore/internal/cli/migrate"
	"bookstore/internal/cli/server"
)

func init() {
	rootCmd.AddCommand(migrate.Command)
	rootCmd.AddCommand(server.Command)
}

var rootCmd = &cobra.Command{
	Use:   "bs",
	Short: "Bookstore application",
}

func Execute() {
	rootCmd.Execute()
}
