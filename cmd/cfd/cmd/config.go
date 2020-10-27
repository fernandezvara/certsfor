package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the bootstrap command
var configCmd = &cobra.Command{
	Use:   "configfile",
	Short: "Creates/updates the configuration file with the current defined values.",
	Long: `Creates the configuration file on the selected location, defaults to (~/.cfd/config.yaml). 
	
	This allows to customize the service before its first run or updates the configuration file if a new version is installed.`,
	Run: configFunc,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configFunc(cmd *cobra.Command, args []string) {

	// create directory
	fmt.Println(viper.GetString(configDBConnectionString))
	er(viper.SafeWriteConfigAs(global.cfgFile))

}
