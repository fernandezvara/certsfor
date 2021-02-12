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
	"time"

	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// listCertificateCmd represents the bootstrap command
var listCertificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert", "certificates"},
	Short:   "List certificates associated with the CA",
	Long: `List certificates associated with the CA

Every CA is identified by a UUID.

The CA ID is required on each request.`,
	Run: listCertificateFunc,
}

func init() {
	listCmd.AddCommand(listCertificateCmd)
	listCertificateCmd.Flags().BoolVar(&global.bool1, "md", false, "Return as markdown formatted text")
	listCertificateCmd.Flags().BoolVar(&global.bool2, "csv", false, "Return as CSV")
}

func listCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv          *service.Service
		collection   string
		certificates map[string]client.Certificate
		t            table.Writer
		keys         []string
		err          error
		ctx          context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	certificates, err = srv.CertificateList(ctx, collection, false) // TODO replace all code with the parsed information
	er(err)

	t = table.NewWriter()

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"Common Name", "Distinguished Name", "Expires in"})

	// sort by cn
	keys = make([]string, 0, len(certificates))
	for key := range certificates {
		keys = append(keys, key)
	}

	for _, key := range keys {

		var (
			expires     string
			now         time.Time = time.Now()
			certificate client.Certificate
		)

		certificate = certificates[key]

		certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)
		er(err)

		if now.Before(certificate.X509Certificate.NotAfter) {
			expires = fmt.Sprintf("%d days (%s)", int64(certificate.X509Certificate.NotAfter.Sub(now).Hours())/int64(24), certificate.X509Certificate.NotAfter.Format(timeFormat))
		} else {
			expires = fmt.Sprintf("Expired (%s)", certificate.X509Certificate.NotAfter.Format(timeFormat))
		}

		t.AppendRow(table.Row{
			key,
			certificate.X509Certificate.Subject.ToRDNSequence(),
			expires,
		})
	}

	if global.bool1 {
		t.RenderMarkdown()
		return
	}

	if global.bool2 {
		t.RenderCSV()
		return
	}

	t.Render()

}
