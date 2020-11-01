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
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/viper"
)

func caTemplate() (request client.APICertificateRequest) {

	request.DN.CN = "ca"
	request.Key = client.RSA4096
	request.ExpirationDays = 365
	return

}

func certTemplate() (request client.APICertificateRequest) {

	request.DN.CN = "common name"
	request.Key = client.RSA4096
	request.ExpirationDays = 90
	return

}

func interactiveCertificate(request *client.APICertificateRequest) {

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
		client.RSA2048,
		client.RSA3072,
		client.RSA4096,
		client.ECDSA224,
		client.ECDSA256,
		client.ECDSA384,
		client.ECDSA521,
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

func saveFiles(ca *manager.CA, bytesCert, bytesKey []byte) {

	if global.certFile != "" {
		er(ioutil.WriteFile(global.certFile, bytesCert, 0400))
	}

	if global.keyFile != "" {
		er(ioutil.WriteFile(global.keyFile, bytesKey, 0400))
	}

	if global.bundleFile != "" {
		er(ioutil.WriteFile(global.keyFile, append(bytesCert, ca.CACertificate()...), 0400))
	}

}
