package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd is the main command for administration on a veil server
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create commands.",
	Long:  `Create commands.`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
