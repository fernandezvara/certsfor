package cmd

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/fernandezvara/certsfor/pkg/manager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// createCertificateCmd represents the bootstrap command
var createCertificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert"},
	Short:   "Creates a new certificate/key pair.",
	Long: `Creates a new certificate/key pair.
	
Every CA is identified by a UUID. 

To operate with the new created CA, you must be pass this ID on each request.`,
	Run: createCertificateFunc,
}

func init() {
	createCmd.AddCommand(createCertificateCmd)
	createCertificateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	createCertificateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location. NOTE: Do not share this file.")
	createCertificateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "File with the answers in YAML format.")
	createCertificateCmd.Flags().StringVar(&global.collection, "ca-id", "", "CA Identifier. (required). [$CFD_CA_ID]")
}

func createCertificateFunc(cmd *cobra.Command, args []string) {

	var (
		srv        *service.Service
		request    structs.APICertificateRequest
		collection string
		ca         *manager.CA
		bytesCert  []byte
		bytesKey   []byte
		bytesInput []byte
		err        error
		ctx        context.Context = context.Background()
	)

	srv = buildService()
	defer srv.Close()

	collection = collectionOrExit()

	ca, err = srv.CAGet(collection)
	fmt.Printf("%v\n", ca)
	er(err)

	// fill from YAML
	// TODO: fill from CSR
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

	bytesCert, bytesKey, err = srv.CertificateSet(ctx, ca, collection, request)
	er(err)

	if global.certFile != "" {
		er(ioutil.WriteFile(global.certFile, bytesCert, 0400))
	}

	if global.keyFile != "" {
		er(ioutil.WriteFile(global.keyFile, bytesKey, 0400))
	}

	fmt.Println("\n\nCertificate Created.")

}

// puedo hablar porque me dice que esta mal porque esta en ingles
