package api

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/fernandezvara/scheduler"
)

// API is the struct that manages the api
type API struct {
	srv     *service.Service
	version string
	server  *rest.REST
	logger  *rest.Logging
	stop    chan os.Signal
}

// New returns the API struct
func New(srv *service.Service, version string) *API {

	return &API{
		srv:     srv,
		version: version,
	}

}

// Start the API
func (a *API) Start(apiPort string, tlsCertificate, tlsKey, tlsCaCert string, remaining int, requireClientCertificate bool, outputPaths, errorOutputPaths []string, debug bool) error {

	var (
		routes            map[string]map[string]rest.APIEndpoint
		cert, key, cacert []byte
		startScheduler    bool
		err               error
	)

	routes = map[string]map[string]rest.APIEndpoint{
		"GET": {
			"/status": {
				Handler: a.getStatus,
				Matcher: []string{""},
			},
			"/v1/ca/:caid/certificates/:cn": {
				Handler: a.getCertificate,
				Matcher: []string{"", "", "", "", "[a-zA-Z0-9.-_]+"},
			},
			"/v1/ca/:caid/certificates": {
				Handler: a.getCertificates,
				Matcher: []string{"", "", "", ""},
			},
		},
		"POST": {
			"/v1/ca": {
				Handler: a.postCA,
				Matcher: []string{"", ""},
			},
		},
		"PUT": {
			"/v1/ca/:caid/certificates/:cn": {
				Handler: a.putCertificate,
				Matcher: []string{"", "", "", "", "[a-zA-Z0-9.-_]+"},
			},
		},
		"DELETE": {
			"/v1/ca/:caid/certificates/:cn": {
				Handler: a.deleteCertificate,
				Matcher: []string{"", "", "", "", "[a-zA-Z0-9.-_]+"},
			},
		},
	}

	a.logger, err = rest.NewLogging(outputPaths, errorOutputPaths, debug)
	if err != nil {
		return err
	}

	// load certificate
	cert, key, cacert, startScheduler, err = getCertificates(tlsCertificate, tlsKey, tlsCaCert, remaining, a.srv)

	a.server, err = rest.New(apiPort, cert, key, cacert, requireClientCertificate, a.logger)
	if err != nil {
		return err
	}

	if startScheduler {

		// look if an certificate update must be done daily
		fn := func() {
			cert, key, _, _, err = getCertificates(tlsCertificate, tlsKey, tlsCaCert, remaining, a.srv)
			a.server.SetCertificate(cert, key)
		}

		scheduler.Every(1).Minutes().NotImmediately().Run(fn)

	}

	a.server.SetupRouter(routes)

	// graceful
	err = a.server.Start()

	return err

}

// Stop the API
func (a *API) Stop() error {

	var err error

	a.logger.Info("api", "initializing server shutdown")

	// close data service
	err = a.srv.Close()
	if err != nil {
		a.server.Shutdown() // try to close it anyway
		return err
	}

	return a.server.Shutdown()

}

func getCertificates(tlsCertificate, tlsKey, tlsCaCert string, remaining int, srv *service.Service) (cert, key, cacert []byte, startScheduler bool, err error) {

	if isUUID(tlsCaCert) {
		// get cert from DB
		var (
			ca  *manager.CA
			crt client.Certificate
		)

		if ca, err = srv.CAGet(tlsCaCert); err != nil {
			return
		}

		if crt, err = srv.CertificateGet(context.Background(), tlsCaCert, tlsCertificate, remaining, false); err != nil {
			return
		}

		cacert = ca.CACertificateBytes()
		cert = crt.Certificate
		key = crt.Key
		startScheduler = true

	} else {
		// ca certificate is a file?
		if cacert, err = fileBytes(tlsCaCert); err != nil {
			return
		}

		if cert, err = fileBytes(tlsCertificate); err != nil {
			return
		}

		if key, err = fileBytes(tlsKey); err != nil {
			return
		}

	}

	return

}

func fileBytes(filename string) ([]byte, error) {

	if filename == "" {
		return []byte{}, nil
	}

	return ioutil.ReadFile(filename)

}

// cfae8b38-57dd-4322-a83f-bc5730689198
func isUUID(uuid string) bool {

	return regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$").MatchString(uuid)

}
