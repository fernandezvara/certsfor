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
	"io/ioutil"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// createCaCmd represents the bootstrap command
var createCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "Creates a new CA certificate/key pair.",
	Long: `Creates a new CA certificate/key pair.
	
Every CA is identified by a UUID. 

To operate with the new created CA, you must be pass this ID on each request.`,
	Run: createCaFunc,
}

func init() {
	createCmd.AddCommand(createCaCmd)
	createCaCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	createCaCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location. NOTE: Do not share this file.")
	createCaCmd.Flags().StringVarP(&global.filename, "file", "f", "", "File with the answers in YAML format.")
}

func createCaFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		request    client.APICertificateRequest
		id         string
		bytesCert  []byte
		bytesKey   []byte
		bytesInput []byte
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	// fill from YAML
	// TODO: fill from CSR
	if global.filename != "" {
		bytesInput, err = ioutil.ReadFile(global.filename)
		er(err)

		err = yaml.Unmarshal(bytesInput, &request)
		er(err)

	} else {
		// interactive: fill Certificate data
		request = caTemplate()
		interactiveCertificate(&request)
	}

	_, id, bytesCert, bytesKey, err = srv.CACreate(ctx, request)

	er(err)

	saveFiles([]byte{}, bytesCert, bytesKey)

	echo(fmt.Sprintf("\n\nCA Created. ID: '%s'\n", id))

}
