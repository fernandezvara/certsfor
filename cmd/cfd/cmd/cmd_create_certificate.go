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
	"io/ioutil"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// createCertificateCmd represents the bootstrap command
var createCertificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert"},
	Short:   "Creates a new certificate/key pair.",
	Long: `Creates a new certificate/key pair.
	
Every CA is identified by a UUID. To operate with the CA, you must be pass this ID on each request.

To write to the standard output (console) the file contents (cert, key, bundle, ca-cert) use 'out' or 'stdout'.`,
	Run: createCertificateFunc,
}

func init() {
	createCmd.AddCommand(createCertificateCmd)
	createCertificateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	createCertificateCmd.Flags().StringVar(&global.caCertFile, "ca-cert", "", "CA Certificate file location.")
	createCertificateCmd.Flags().StringVarP(&global.bundleFile, "bundle", "b", "", "Bundle file location.")
	createCertificateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location. NOTE: Do not share this file.")
	createCertificateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "File with the answers in YAML format.")
	createCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
}

func createCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		request    client.APICertificateRequest
		collection string
		bytesCA    []byte
		bytesCert  []byte
		bytesKey   []byte
		bytesInput []byte
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	// TODO: fill from CSR

	// fill from YAML
	if global.filename != "" {
		bytesInput, err = ioutil.ReadFile(global.filename)
		er(err)

		err = yaml.Unmarshal(bytesInput, &request)
		er(err)

	} else {
		// interactive: fill Certificate data
		request = certTemplate()
		interactiveCertificate(&request)
	}

	bytesCA, bytesCert, bytesKey, err = srv.CertificateSet(ctx, collection, request)
	er(err)

	saveFiles(bytesCA, bytesCert, bytesKey)

	echo("\n\nCertificate Created.")

}
