package certinfo

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

// errors
var (
	ErrNoTLSConnection = errors.New("URL is not a TLS connection")
	ErrUnparseableFile = errors.New("unparseable file")
)

const (
	timeFormat = "02/01/2006 15:04"
)

// NewFromURL returns a CertInfo struct for a URL
func NewFromURL(uri string, timeout int64, insecure bool) (certInfo CertInfo, err error) {

	var (
		ipPort string
		conn   *tls.Conn
		dialer *net.Dialer
	)

	certInfo = CertInfo{}
	certInfo.TimeFormat = timeFormat

	ipPort, err = cleanURI(uri)
	if err != nil {
		return
	}

	dialer = &net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	conn, err = tls.DialWithDialer(dialer, "tcp", ipPort, &tls.Config{
		InsecureSkipVerify: insecure,
	})

	switch err.(type) {
	case tls.RecordHeaderError:
		err = ErrNoTLSConnection
		return
	case nil:
		// continue
	default:
		return
	}
	defer conn.Close()

	certInfo.certs = conn.ConnectionState().PeerCertificates

	return

}

// NewFromBytes returns a CertInfo from a certificate bytes
func NewFromBytes(certBytes []byte) (certInfo CertInfo, err error) {

	certInfo = CertInfo{}
	certInfo.TimeFormat = timeFormat

	err = certInfo.fromBytes(certBytes)
	return
}

// NewFromFile returns a CertInfo from a certificate file
func NewFromFile(filename string) (certinfo CertInfo, err error) {

	var certBytes []byte

	certBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	return NewFromBytes(certBytes)

}

func cleanURI(input string) (ipPort string, err error) {

	var (
		u *url.URL
	)

	if !strings.Contains(input, "://") {
		input = fmt.Sprintf("https://%s", input)
	}

	u, err = url.Parse(input)
	if err != nil {
		return
	}

	ipPort = u.Host
	if u.Port() == "" {
		ipPort = fmt.Sprintf("%s:443", u.Host)
	}

	return
}

// CertInfo is the struct used to query and return the certificate/s information
type CertInfo struct {
	certs      []*x509.Certificate
	TimeFormat string // time format to use
}

func (c *CertInfo) fromBytes(certBytes []byte) (err error) {

	var (
		cert *x509.Certificate
	)

	block, _ := pem.Decode([]byte(certBytes))
	if block == nil {
		return ErrUnparseableFile
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}
	c.certs = append(c.certs, cert)

	return
}

// Show returns the information of the certificates formatted as table
func (c *CertInfo) Show(markdown, csv bool) {

	var (
		t table.Writer
	)

	t = table.NewWriter()

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"Common Name", "Distinguished Name", "Not Before", "Expires", "SANs", "CA?"})

	for _, cert := range c.certs {

		var (
			now     time.Time = time.Now()
			expires string
		)
		if now.Before(cert.NotAfter) {
			expires = fmt.Sprintf("%d days (%s)", int64(cert.NotAfter.Sub(now).Hours())/int64(24), cert.NotAfter.Format(c.TimeFormat))
		} else {
			expires = fmt.Sprintf("Expired (%s)", cert.NotAfter.Format(c.TimeFormat))
		}

		t.AppendRow(table.Row{
			cert.Subject.CommonName,
			cert.Subject.ToRDNSequence(),
			cert.NotBefore.In(time.UTC).Format(c.TimeFormat),
			expires,
			strings.Join(cert.DNSNames, "\n"),
			cert.IsCA,
		})

	}

	if markdown {
		t.RenderMarkdown()
		return
	}

	if csv {
		t.RenderCSV()
		return
	}

	t.Render()

}

// Certificates returns the array of certificates
func (c *CertInfo) Certificates() []*x509.Certificate {

	return c.certs

}
