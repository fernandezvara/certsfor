package cmd

import (
	"os"
	"path/filepath"

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
	if _, err := os.Stat(filepath.Dir(global.cfgFile)); os.IsNotExist(err) {
		er(os.MkdirAll(filepath.Dir(global.cfgFile), 0750))
	}

	er(viper.WriteConfigAs(global.cfgFile))

}
