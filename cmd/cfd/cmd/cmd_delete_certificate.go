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

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/spf13/cobra"
)

// deleteCertificateCmd represents the bootstrap command
var deleteCertificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert"},
	Short:   "Deletes a certificate by its Common Name from the store",
	Long: `Deletes a certificate by its Common Name from the store
	
Every CA is identified by a UUID. 

The CA ID is required on each request.`,
	Run: deleteCertificateFunc,
}

func init() {
	deleteCmd.AddCommand(deleteCertificateCmd)
	deleteCertificateCmd.Flags().BoolVarP(&global.bool1, "yes", "y", false, "Asumme yes to the prompts (is assumed if --quiet)")
	deleteCertificateCmd.Flags().StringVar(&global.cn, "cn", "", "Common Name. (required)")
	deleteCertificateCmd.MarkFlagRequired("cn")
}

func deleteCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		collection string
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	// if quiet assume yes
	if global.quiet {
		global.bool1 = global.quiet
	}

	// run interactively?
	if !global.bool1 {
		global.bool1, err = promptTrueFalseBool("Are you sure?", "Yes", "No", false)
		er(err)
	}

	if global.bool1 {
		_, err = srv.CertificateDelete(ctx, collection, global.cn)
		er(err)

		echo("\n\nCertificate Deleted.")
	} else {
		echo("\n\nOperation Cancelled.")
	}

}
