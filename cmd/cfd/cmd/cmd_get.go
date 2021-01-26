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

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
)

// getCertificateCmd represents the bootstrap command
var getCertificateCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a certificate/key from the store.",
	Long: `Get a certificate/key from the store.
	
To write to the standard output (console) the file contents (cert, key, bundle, ca-cert) use 'out' or 'stdout'.`,
	Run: getCertificateFunc,
}

func init() {
	rootCmd.AddCommand(getCertificateCmd)
	getCertificateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	getCertificateCmd.Flags().StringVar(&global.caCertFile, "ca-cert", "", "CA Certificate file location.")
	getCertificateCmd.Flags().StringVarP(&global.bundleFile, "bundle", "b", "", "Bundle file location.")
	getCertificateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location.")
	getCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
	getCertificateCmd.Flags().StringVar(&global.cn, "cn", "", "Common Name. (required).")
	getCertificateCmd.Flags().IntVar(&global.remaining, "renew", 20, "Time (expresed as percent) to be used to determine if the certificate must be renewed (defaults to 20 %). Key remains the same.")
	getCertificateCmd.MarkFlagRequired("cn")
}

func getCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		collection string
		cert       client.Certificate
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	cert, err = srv.CertificateGet(ctx, collection, global.cn, global.remaining)
	er(err)

	saveFiles(cert.CACertificate, cert.Certificate, cert.Key)

	fmt.Println("\n\nCertificate Retrieved.")

}
