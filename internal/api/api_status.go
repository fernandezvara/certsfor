package api

import (
	"net/http"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// GetStatus GET /status
func (a *API) getStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		response client.APIStatus
	)

	rest.Response(w, response, nil, 200, "")

}
