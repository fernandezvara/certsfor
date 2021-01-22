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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

// webserverCmd represents the bootstrap command
var webserverCmd = &cobra.Command{
	Use:   "webserver",
	Short: "Starts a simple webserver using a configured certificate.",
	Long:  `Starts a simple webserver using a configured certificate.`,
	Run:   webserverFunc,
}

func init() {
	startCmd.AddCommand(webserverCmd)
	webserverCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
	webserverCmd.Flags().StringVar(&global.cn, "cn", "", "Common Name. (required).")
	webserverCmd.Flags().IntVar(&global.remaining, "renew", 20, "Time (expresed as percent) to be used to determine if the certificate must be renewed (defaults to 20 %). Key remains the same.")
	webserverCmd.Flags().StringVar(&global.listen, "listen", "0.0.0.0:8443", "IP:TCP Port where the server will be served. Defaults to all network interfaces and port 8443.")
	webserverCmd.Flags().StringVar(&global.root, "root", ".", "Directory where the files reside, defaults to current (.).")
	webserverCmd.MarkFlagRequired("cn")
}

func webserverFunc(cmd *cobra.Command, args []string) {

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

	// start web server
	var (
		httpServer  http.Server
		router      *httprouter.Router
		fullPath    string
		tlsConfig   tls.Config
		certificate tls.Certificate
	)

	fullPath = path.Clean(global.root)

	router = httprouter.New()
	router.ServeFiles("/*filepath", http.Dir(fullPath))

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert.CACertificate)

	tlsConfig.ClientCAs = caCertPool
	tlsConfig.BuildNameToCertificate()

	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}
	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}

	tlsConfig.NextProtos = []string{"h2", "http/1.1"}

	certificate, err = tls.X509KeyPair(cert.Certificate, cert.Key)
	er(err)

	tlsConfig.Certificates = []tls.Certificate{certificate}

	httpServer.Addr = global.listen
	httpServer.Handler = router
	httpServer.TLSConfig = &tlsConfig

	// graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err = httpServer.ListenAndServeTLS("", ""); err != nil {
			if err != http.ErrServerClosed {
				fmt.Println("Error on WebServer", err.Error())
				os.Exit(1)
			}
		}
	}()

	<-stop

	httpServer.Shutdown(context.Background())

}
