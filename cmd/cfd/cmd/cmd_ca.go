package cmd

import (
	"github.com/spf13/cobra"
)

// caCmd is the main command for administration on a veil server
var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "CA administration commands.",
	Long:  `CA administration commands.`,
}

func init() {
	rootCmd.AddCommand(caCmd)
}
