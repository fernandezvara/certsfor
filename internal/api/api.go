package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/rest"
)

// API is the struct that manages the api
type API struct {
	srv     *service.Service
	version string
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
		server *rest.REST
		logger *rest.Logging
		err    error
	)

	routes = map[string]map[string]rest.APIEndpoint{
		"GET": {
			"/status": {
				Handler: a.getStatus,
				Matcher: []string{""},
			},
		},
	}

	logger, err = rest.NewLogging(outputPaths, errorOutputPaths, debug)
	if err != nil {
		return err
	}

	server, err = rest.New(apiPort, tlsCertificate, tlsKey, tlsCACert, logger)
	if err != nil {
		return err
	}

	server.SetupRouter(routes)

	// graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err = server.Start(); err != nil {
			if err != http.ErrServerClosed {
				fmt.Println("Error on API.Server", err.Error())
				os.Exit(1)
			}
		}
	}()

	<-stop

	logger.Info("api", "initializing server shutdown")

	// close data service
	err = a.srv.Close()
	if err != nil {
		server.Shutdown() // try to close it anyway
		return err
	}

	return server.Shutdown()

}
