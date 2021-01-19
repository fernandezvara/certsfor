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

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/spf13/cobra"
)

// getStatusCmd represents the bootstrap command
var getStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of local / remote API.",
	Long:  `Status of local / remote API.`,
	Run:   getStatusFunc,
}

func init() {
	rootCmd.AddCommand(getStatusCmd)
	getStatusCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	getStatusCmd.Flags().StringVarP(&global.bundleFile, "bundle", "b", "", "Bundle file location.")
	getStatusCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location.")
}

func getStatusFunc(cmd *cobra.Command, args []string) {

	var (
		srv *service.Service
		err error
	)

	srv = buildService()
	defer srv.Close()

	status, err := srv.Status()

	fmt.Println(status, err)

}
