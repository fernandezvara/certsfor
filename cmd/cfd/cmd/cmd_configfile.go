/*
Copyright © 2020 @fernandezvara

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
