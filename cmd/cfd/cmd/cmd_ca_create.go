package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// caCreateCmd represents the bootstrap command
var caCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new CA certificate/key pair.",
	Long: `Creates a new CA certificate/key pair.
	
Every CA is identified by a UUID. 

To operate with the new created CA, you must be pass this ID on each request.`,
	Run: caCreateFunc,
}

func init() {
	caCmd.AddCommand(caCreateCmd)
	caCreateCmd.Flags().StringVarP(&global.certFile, "cert", "c", "", "Certificate file location.")
	caCreateCmd.Flags().StringVarP(&global.keyFile, "key", "k", "", "Key file location. NOTE: Do not share this file.")
	caCreateCmd.Flags().StringVarP(&global.filename, "file", "f", "", "File with the answers in YAML format.")
}

func caTemplate() (request structs.APICertificateRequest) {

	request.DN.CN = "ca"
	request.Key = structs.RSA3072
	request.ExpirationDays = 365
	return

}

func caCreateFunc(cmd *cobra.Command, args []string) {

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

func interactiveCertificate(request *structs.APICertificateRequest) {

	var (
		err     error
		expires string
	)

	request.DN.CN, err = promptText("Common Name", request.DN.CN, validationRequired)
	er(err)
	request.DN.C, err = promptText("Country (optional)", request.DN.C, nil)
	er(err)
	request.DN.P, err = promptText("Province (optional)", request.DN.P, nil)
	er(err)
	request.DN.L, err = promptText("Locality (optional)", request.DN.L, nil)
	er(err)
	request.DN.PC, err = promptText("Postal Code (optional)", request.DN.PC, nil)
	er(err)
	request.DN.ST, err = promptText("Street (optional)", request.DN.ST, nil)
	er(err)
	request.DN.O, err = promptText("Organization (optional)", request.DN.O, nil)
	er(err)
	request.DN.OU, err = promptText("Organizational Unit (optional)", request.DN.OU, nil)
	er(err)
	expires, err = promptText("Expires in (days)", strconv.Itoa(int(request.ExpirationDays)), nil)
	er(err)
	request.ExpirationDays, err = strconv.ParseInt(expires, 10, 64)
	er(err)

	items := []string{
		"RSA (2048 bytes)",
		"RSA (3072 bytes)",
		"RSA (4096 bytes)",
		"ECDSA (EC-224)",
		"ECDSA (EC-256)",
		"ECDSA (EC-384)",
		"ECDSA (EC-521)",
	}
	values := []string{
		structs.RSA2048,
		structs.RSA3072,
		structs.RSA4096,
		structs.ECDSA224,
		structs.ECDSA256,
		structs.ECDSA384,
		structs.ECDSA521,
	}

	request.Key, err = promptSelection("Key Algorithm", items, values, 2)
	er(err)

}
