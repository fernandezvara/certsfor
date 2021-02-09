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
	"fmt"
	"os"

	"github.com/fernandezvara/certsfor/internal/certinfo"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
)

// infoCertificateCmd represents the bootstrap command
var infoCertificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert"},
	Short:   "Show (friendly) information about a certificate.",
	Long: `Show information about a certificate. Allows to get data from:
 
  - File
  - Common Name
  - Remote URL

`,
	Run: infoCertificateFunc,
}

var (
	timeout  int64
	insecure bool
)

func init() {
	infoCmd.AddCommand(infoCertificateCmd)
	infoCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required if --cn) [$CFD_CA_ID]")
	infoCertificateCmd.Flags().StringVar(&global.cn, "cn", "", "Common name to query for information.")
	infoCertificateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "Source file with the certificate.")
	infoCertificateCmd.Flags().StringVarP(&global.listen, "url", "u", "", "URL to get information from.")
	infoCertificateCmd.Flags().Int64Var(&timeout, "timeout", 5, "Timeout for network calls (only used if --url is especified).")
	infoCertificateCmd.Flags().BoolVar(&insecure, "insecure", false, "Insecure allows skip verification of the server certificate (only used if --url is especified).")
	infoCertificateCmd.Flags().BoolVar(&global.bool1, "markdown", false, "Return the data formatted as markdown.")
	infoCertificateCmd.Flags().BoolVar(&global.bool2, "csv", false, "Show as CSV.")
}

func infoCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		certInfo   certinfo.CertInfo
		response   client.Certificate
		collection string
		err        error
		ctx        context.Context = context.Background()
	)

	// if cn ensure ca-id getInfo from certificate
	if global.cn != "" {

		collection = collectionOrExit()
		srv = buildService()
		defer srv.Close()

		response, err = srv.CertificateGet(ctx, collection, global.cn, 0)
		er(err)

		certInfo, err = certinfo.NewFromBytes(response.Certificate)
		er(err)

	} else {
		// if file getInfo.fromfile()
		if global.filename != "" {

			certInfo, err = certinfo.NewFromFile(global.filename)
			er(err)

		} else {

			if global.listen == "" {
				fmt.Println("Select the certificate to get its information")
				os.Exit(1)
			}

			// if url getInfo.from URL()
			certInfo, err = certinfo.NewFromURL(global.listen, timeout, insecure)
			er(err)

		}

	}

	certInfo.Show(global.bool1, global.bool2)

}
