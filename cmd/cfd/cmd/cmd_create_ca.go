package cmd

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/structs"
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
		request    structs.APICertificateRequest
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

	if global.certFile != "" {
		er(ioutil.WriteFile(global.certFile, bytesCert, 0400))
	}

	if global.keyFile != "" {
		er(ioutil.WriteFile(global.keyFile, bytesKey, 0400))
	}

	fmt.Printf("\n\nCA Created. ID: '%s'\n", id)

}

// puedo hablar porque me dice que esta mal porque esta en ingles
