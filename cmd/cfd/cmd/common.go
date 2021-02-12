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
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/viper"
	"software.sslmate.com/src/go-pkcs12"
)

func caTemplate() (request client.APICertificateRequest) {

	request.DN.CN = "ca"
	request.Key = client.RSA4096
	request.ExpirationDays = 365
	return

}

func certTemplate() (request client.APICertificateRequest) {

	request.DN.CN = "common-name"
	request.Key = client.RSA4096
	request.ExpirationDays = 90
	return

}

func interactiveCertificate(request *client.APICertificateRequest, isCertificate bool) {

	var (
		err     error
		expires string
	)

	request.DN.CN, err = promptText("Common Name", request.DN.CN, validationRequired)
	er(err)
	request.DN.C, err = promptText("Country (optional)", request.DN.C, nil)
	er(err)
	request.DN.P, err = promptText("Province / State (optional)", request.DN.P, nil)
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

	// do no ask for CA
	if isCertificate {
		request.Client, err = promptTrueFalseBool("Client Certificate?", "Yes", "No", false)
		er(err)
	}

}

func collectionOrExit() (collection string) {

	collection = viper.GetString("ca-id")
	if global.collection != "" {
		collection = global.collection
	}

	if collection == "" {
		echo("\n\nCA ID missing. Cannot continue.")
		os.Exit(1)
	}

	return
}

func saveFiles(ca, bytesCert, bytesKey []byte) {

	var (
		pfxData                          []byte
		blockCert, blockCaCert, blockKey *pem.Block
		cert, caCert                     *x509.Certificate
		key                              interface{}
		err                              error
	)

	// pkcs12
	if global.pfxFile != "" {

		blockCert, _ = pem.Decode(bytesCert)
		cert, err = x509.ParseCertificate(blockCert.Bytes)
		er(err)

		blockCaCert, _ = pem.Decode(ca)
		caCert, err = x509.ParseCertificate(blockCaCert.Bytes)
		er(err)

		blockKey, _ = pem.Decode(bytesKey)
		key, err = x509.ParsePKCS8PrivateKey(blockKey.Bytes)
		er(err)

		pfxData, err = pkcs12.Encode(rand.Reader, key, cert, []*x509.Certificate{caCert}, global.pfxPassword)
		er(err)

		saveOrShowFile(global.pfxFile, pfxData, 0400)

	}

	saveOrShowFile(global.certFile, bytesCert, 0400)
	saveOrShowFile(global.keyFile, bytesKey, 0400)
	saveOrShowFile(global.bundleFile, append(bytesCert, ca...), 0400)
	saveOrShowFile(global.caCertFile, ca, 0400)

}

func saveOrShowFile(file string, contents []byte, perm os.FileMode) {

	switch file {
	case "out", "stdout":
		fmt.Print(string(contents))
	case "":
		// do nothing
	default:
		er(ioutil.WriteFile(file, contents, perm))
	}

}

func echo(message string) {
	if !global.quiet {
		fmt.Println(message)
	}
}
