package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/spf13/viper"
)

func caTemplate() (request structs.APICertificateRequest) {

	request.DN.CN = "ca"
	request.Key = structs.RSA4096
	request.ExpirationDays = 365
	return

}

func certTemplate() (request structs.APICertificateRequest) {

	request.DN.CN = "common name"
	request.Key = structs.RSA4096
	request.ExpirationDays = 365
	return

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
	request.SAN, err = promptArray("Hosts and IPs. (blank if finish)")
	er(err)
	expires, err = promptText("Expires in (days)", strconv.Itoa(int(request.ExpirationDays)), validationInteger)
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

func collectionOrExit() (collection string) {

	collection = viper.GetString("ca-id")
	if global.collection != "" {
		collection = global.collection
	}

	if collection == "" {
		fmt.Println("\n\nCA ID missing. Cannot continue.")
		os.Exit(1)
	}

	return
}
