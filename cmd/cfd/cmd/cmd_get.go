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

	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCertificateCmd represents the bootstrap command
var getCertificateCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a certificate/key from the store.",
	Long:  `Get a certificate/key from the store.`,
	Run:   getCertificateFunc,
}

func init() {
	rootCmd.AddCommand(getCertificateCmd)
	getCertificateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	getCertificateCmd.Flags().StringVarP(&global.bundleFile, "bundle", "b", "", "Bundle file location.")
	getCertificateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location.")
	getCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
	getCertificateCmd.Flags().StringVar(&global.collection, "cn", "", "Common Name. (required).")
	getCertificateCmd.MarkFlagRequired("cn")
}

func getCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		collection string
		ca         *manager.CA
		cert       client.Certificate
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	ca, err = srv.CAGet(collection)
	er(err)

	cert, err = srv.CertificateGet(ctx, collection, viper.GetString("cn"))
	er(err)

	saveFiles(ca, cert.Certificate, cert.Key)

	fmt.Println("\n\nCertificate Created.")

}
