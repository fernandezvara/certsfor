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
	"context"
	"os"
	"os/signal"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/api"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the bootstrap command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts as REST API service allowing remote cfd clients.",
	Long:  `Starts as REST API service allowing remote cfd clients.`,
	Run:   apiFunc,
}

func init() {
	startCmd.AddCommand(apiCmd)
}

func apiFunc(cmd *cobra.Command, args []string) {

	var (
		sto store.Store
		srv *service.Service
		a   *api.API
		err error
	)

	sto, err = store.Open(context.Background(), viper.GetString(configDBType), viper.GetString(configDBConnectionString))
	er(err)
	srv = service.NewAsServer(sto, Version)

	a = api.New(srv, Version)
	go a.Start(viper.GetString(configAPIAddr),
		viper.GetString(configTLSCertificate),
		viper.GetString(configTLSKey),
		viper.GetString(configTLSCA),
		global.remaining,
		viper.GetBool(configTLSRequireClientCertificate),
		viper.GetStringSlice(configAPIAccessLog),
		viper.GetStringSlice(configAPIErrorLog),
		viper.GetBool(configAPIDebugLog),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	echo("\nStopping API...\n")

	err = a.Stop()
	er(err)

}
