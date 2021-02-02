package api

import (
	"os"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/rest"
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
func (a *API) Start(apiPort string, tlsCertificate, tlsKey, tlsCACert []byte, outputPaths, errorOutputPaths []string, debug bool) error {

	var (
		routes map[string]map[string]rest.APIEndpoint
		err    error
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
	}

	a.logger, err = rest.NewLogging(outputPaths, errorOutputPaths, debug)
	if err != nil {
		return err
	}

	a.server, err = rest.New(apiPort, tlsCertificate, tlsKey, tlsCACert, a.logger)
	if err != nil {
		return err
	}

	a.server.SetupRouter(routes)

	// graceful
	return a.server.Start()

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
