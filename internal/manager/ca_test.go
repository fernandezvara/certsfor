package manager

import (
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestNewCA(t *testing.T) {

	var request client.APICertificateRequest

	request.DN.CN = "this is a test"
	request.DN.C = "ES"
	request.DN.L = "locality"
	request.DN.O = "organization"
	request.DN.OU = "ourganization unit"
	request.DN.P = "province"
	request.DN.PC = "postalCode"
	request.DN.ST = "street"

	request.ExpirationDays = 90
	request.Key = client.RSA4096

	certFile, keyFile, err := New(request)

	assert.Nil(t, err)
	assert.Greater(t, len(certFile), 0)
	assert.Greater(t, len(keyFile), 0)

	newCA2, err := FromBytes(certFile, keyFile)
	assert.Nil(t, err)
	assert.NotNil(t, newCA2)

	assert.Len(t, newCA2.ca.Issuer.SerialNumber, 0)

	assert.True(t, newCA2.ca.IsCA)

	assert.Equal(t, newCA2.ca.Subject.CommonName, request.DN.CN)

	assert.Len(t, newCA2.ca.Subject.Organization, 1)
	assert.Len(t, newCA2.ca.Subject.Country, 1)
	assert.Len(t, newCA2.ca.Subject.Province, 1)
	assert.Len(t, newCA2.ca.Subject.Locality, 1)
	assert.Len(t, newCA2.ca.Subject.StreetAddress, 1)
	assert.Len(t, newCA2.ca.Subject.PostalCode, 1)

	assert.Equal(t, newCA2.ca.Subject.Organization[0], request.DN.O)
	assert.Equal(t, newCA2.ca.Subject.Country[0], request.DN.C)
	assert.Equal(t, newCA2.ca.Subject.Province[0], request.DN.P)
	assert.Equal(t, newCA2.ca.Subject.Locality[0], request.DN.L)
	assert.Equal(t, newCA2.ca.Subject.PostalCode[0], request.DN.PC)

}

func TestCertificateWithoutCommonName(t *testing.T) {

	request := client.APICertificateRequest{
		ExpirationDays: 90,
		Key:            client.RSA4096,
	}

	certFile, keyFile, err := New(request)
	assert.Error(t, err, ErrCommonNameBlank)
	assert.Len(t, certFile, 0)
	assert.Len(t, keyFile, 0)
}

func TestDifferentKeyTypes(t *testing.T) {

	var request client.APICertificateRequest

	for _, typ := range []string{client.RSA2048, client.RSA3072, client.RSA4096, client.ECDSA224, client.ECDSA256, client.ECDSA384, client.ECDSA521} {

		request.DN.CN = "Test"
		request.ExpirationDays = 90
		request.Key = typ

		certFile, keyFile, err := New(request)

		fmt.Println(string(certFile))
		assert.Nil(t, err)
		assert.Greater(t, len(certFile), 0)
		assert.Greater(t, len(keyFile), 0)

	}

	request.DN.CN = "Test"
	request.DN.ST = "street"
	request.ExpirationDays = 90
	request.Key = "invalid"

	certFile, keyFile, err := New(request)

	assert.Error(t, ErrKeyInvalid, err)
	assert.Len(t, certFile, 0)
	assert.Len(t, keyFile, 0)

}

func TestSANsAndAPIrequests(t *testing.T) {

	var request client.APICertificateRequest

	request.DN.CN = "Test"
	request.ExpirationDays = 90
	request.Key = client.ECDSA521

	ip1, ip2 := "192.168.1.1", "123.123.123.123"
	dns1, dns2 := "*.iswildcard.com", "www.example.com"
	uri1, uri2 := "https://www.example1.com", "http://isuri.isuri.com"
	request.SAN = []string{ip1, ip2, dns1, dns2, uri1, uri2}

	certFile, keyFile, err := New(request)

	assert.Nil(t, err)
	assert.Greater(t, len(certFile), 0)
	assert.Greater(t, len(keyFile), 0)

	otherManager, err := FromBytes(certFile, keyFile)
	assert.Nil(t, err)

	assert.Len(t, otherManager.ca.URIs, 2)
	u1, err := url.Parse(uri1)
	assert.Nil(t, err)
	assert.Contains(t, otherManager.ca.URIs, u1)
	u2, err := url.Parse(uri2)
	assert.Nil(t, err)
	assert.Contains(t, otherManager.ca.URIs, u2)

	assert.Len(t, otherManager.ca.IPAddresses, 2)
	assert.Equal(t, otherManager.ca.IPAddresses[0].String(), net.ParseIP(ip1).String())
	assert.Equal(t, otherManager.ca.IPAddresses[1].String(), net.ParseIP(ip2).String())

	assert.Len(t, otherManager.ca.DNSNames, 2)
	assert.Contains(t, otherManager.ca.DNSNames, dns1)
	assert.Contains(t, otherManager.ca.DNSNames, dns2)

	assert.Equal(t, certFile, otherManager.CACertificateBytes())
	assert.Equal(t, request.DN.CN, otherManager.CACertificate().Subject.CommonName)

	var newRequest client.APICertificateRequest

	newRequest.DN.CN = "Client"
	newRequest.ExpirationDays = 90
	newRequest.Key = client.ECDSA521
	newRequest.Client = true

	certFile, keyFile, err = otherManager.CreateCertificateFromAPI(newRequest)
	assert.Nil(t, err)
	assert.Greater(t, len(certFile), 0)
	assert.Greater(t, len(keyFile), 0)

	var newRequest2 client.APICertificateRequest

	newRequest2.DN.CN = "email"
	newRequest2.ExpirationDays = 90
	newRequest2.Key = client.ECDSA521
	newRequest2.SAN = []string{"email@example.com"}

	certFile, keyFile, err = otherManager.CreateCertificateFromAPI(newRequest2)
	assert.Nil(t, err)
	assert.Greater(t, len(certFile), 0)
	assert.Greater(t, len(keyFile), 0)

	var newRequest3 client.APICertificateRequest

	newRequest3.DN.CN = "site"
	newRequest3.ExpirationDays = 90
	newRequest3.Key = client.ECDSA521
	newRequest3.SAN = []string{"www.example.org"}

	certFile, keyFile, err = otherManager.CreateCertificateFromAPI(newRequest3)
	assert.Nil(t, err)
	assert.Greater(t, len(certFile), 0)
	assert.Greater(t, len(keyFile), 0)

	var newRequest4 client.APICertificateRequest

	newRequest4.DN.CN = "with eror"
	newRequest4.ExpirationDays = 90
	newRequest4.Key = "invalid"
	newRequest4.SAN = []string{"www.example.org"}

	certFile, keyFile, err = otherManager.CreateCertificateFromAPI(newRequest4)
	assert.Error(t, err)
	assert.Len(t, certFile, 0)
	assert.Len(t, keyFile, 0)

}

func TestUnparseableFiles(t *testing.T) {

	goodCert := `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

	goodKey := `-----BEGIN PRIVATE KEY-----
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
-----END PRIVATE KEY-----`

	badCert := `-----BEGIN CERTIFICATE-----
MIIF+TCCA+GgAwIBAgIIFj/oIAfys1AwDQYJKoZIhvcNAQELBQAwgYoxCzAJBgNV
BAYTAkVTMREwDwYDVQQIEwhwcm92aW5jZTERMA8GA1UEBxMIbG9jYWxpdHkxFjAU
BgNVBAkTDXN0cmVldGFkZHJlc3MxDjAMBgNVBBETBTEyMzQ1MRcwFQYDVQQKEw5U
-----END CERTIFICATE-----`

	badCert2 := `MIIF+TCCA+GgAwIBAgIIFj/oIAfys1AwDQYJKoZIhvcNAQELBQAwgYoxCzAJBgNV
BAYTAkVTMREwDwYDVQQIEwhwcm92aW5jZTERMA8GA1UEBxMIbG9jYWxpdHkxFjAU
BgNVBAkTDXN0cmVldGFkZHJlc3MxDjAMBgNVBBETBTEyMzQ1MRcwFQYDVQQKEw5U
-----END CERTIFICATE-----`

	badKey := `MIIJKQIBAAKCAgEApfwi3ffMS3sVyp4nt4++oon3hC+1tbHbyy7coaI8qBtT++tF
la06r0aAby+KxMaavqurOBAI+WZSAOhhEmgM2jB1lDHpnO1xTSMYG/4rLJ3yNqnj
3LAfcJHH8OrlSmkn1+sBNjd/nZtpl3n00CmNffhgHMs9AZC6Sqs9mDj35HTl+Okr`

	badKey2 := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----` // flip last 5 characters

	errCA, err := FromBytes([]byte(badCert), []byte(badKey))

	assert.Nil(t, errCA)
	assert.Error(t, err)
	assert.Equal(t, "asn1: syntax error: data truncated", err.Error())

	errCA, err = FromBytes([]byte(goodCert), []byte(badKey))

	assert.Nil(t, errCA)
	assert.Error(t, err)
	assert.IsType(t, err, ErrUnparseableFile)

	errCA, err = FromBytes([]byte(goodCert), []byte(badKey2))

	assert.Nil(t, errCA)
	assert.Error(t, err)
	assert.IsType(t, err, ErrUnparseableFile)

	errCA, err = FromBytes([]byte(badCert), []byte(goodKey))

	assert.Nil(t, errCA)
	assert.Error(t, err)
	assert.Equal(t, "asn1: syntax error: data truncated", err.Error())

	errCA, err = FromBytes([]byte(badCert2), []byte(goodKey))

	assert.Nil(t, errCA)
	assert.Error(t, err)
	assert.IsType(t, err, ErrUnparseableFile)

}
