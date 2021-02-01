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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cfd",
	Short: "certsfor.dev",
	Long: `
'cfd' is a tool for certificate workflow for development environments.
	
Full documentation: https://www.certsfor.dev/
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&global.cfgFile, configFileString, configFileStringDefault, "Configuration file")
	rootCmd.PersistentFlags().BoolVarP(&global.quiet, configQuiet, configQuietShort, configQuietDefault, "Suppress the command output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// default configuration
	// db
	viper.SetDefault(configDBType, configDBTypeDefault)
	viper.BindEnv(configDBType, configDBTypeEnv)
	viper.SetDefault(configDBConnectionString, configDBConnectionStringDefault)
	viper.BindEnv(configDBConnectionString, configDBConnectionStringEnv)

	// api
	viper.SetDefault(configAPIAddr, configAPIAddrDefault)
	viper.BindEnv(configAPIAddr, configAPIAddrEnv)
	viper.SetDefault(configAPIAccessLog, configAPIAccessLogDefault)
	viper.BindEnv(configAPIAccessLog, configAPIAccessLogEnv)
	viper.SetDefault(configAPIErrorLog, configAPIErrorLogDefault)
	viper.BindEnv(configAPIErrorLog, configAPIErrorLogEnv)
	viper.SetDefault(configAPIDebugLog, configAPIDebugLogDefault)
	viper.BindEnv(configAPIDebugLog, configAPIDebugLogEnv)
	// viper.SetDefault(configWebEnabled, configWebEnabledDefault)
	// viper.BindEnv(configWebEnabled, configWebEnabledEnv)

	// api - client
	viper.SetDefault(configAPIEnabled, configAPIEnabledDefault)
	viper.BindEnv(configAPIEnabled, configAPIEnabledEnv)

	// tls
	viper.SetDefault(configTLSCA, configTLSCADefault)
	viper.BindEnv(configTLSCA, configTLSCAEnv)
	viper.SetDefault(configTLSCertificate, configTLSCertificateDefault)
	viper.BindEnv(configTLSCertificate, configTLSCertificateEnv)
	viper.SetDefault(configTLSKey, configTLSKeyDefault)
	viper.BindEnv(configTLSKey, configTLSKeyEnv)
	viper.SetDefault(configTLSUseForce, configTLSUseForceDefault)
	viper.BindEnv(configTLSUseForce, configTLSUseForceEnv)

	// ca id
	viper.SetDefault(configCAID, configCAIDDefault)
	viper.BindEnv(configCAID, configCAIDEnv)

	// quiet mode
	viper.SetDefault(configQuiet, configQuietDefault)
	viper.BindEnv(configQuiet, configQuietEnv)

	global.quiet = viper.GetBool(configQuiet)

	// force quiet mode? when echoing to stdout we must clean the output
	global.quiet = isOutput()

	if os.Getenv("CFD_CONFIG") != "" {
		global.cfgFile = os.Getenv("CFD_CONFIG")
	}

	if global.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(global.cfgFile)
	} else {
		viper.AddConfigPath(pathFromHome(".cfd"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		echo(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}

func isOutput() bool {

	for _, value := range []string{global.certFile, global.keyFile, global.bundleFile, global.caCertFile} {
		if value == "out" || value == "stdout" {
			return true
		}
	}

	return false

}

func er(err error) {
	if err != nil {
		fmt.Printf("err: %T \n%s\n", err, err.Error())
		os.Exit(1)
	}
}
