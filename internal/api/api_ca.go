package api

import (
	"fmt"
	"net/http"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// postCA POST /v1/ca
func (a *API) postCA(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		request  client.APICertificateRequest
		response client.Certificate
		err      error
	)

	err = rest.GetFromBody(r, &request)
	if err != nil {
		rest.BadRequest(w, r, "")
		return
	}

	response.Request = request

	_, response.CAID, response.Certificate, response.Key, err = a.srv.CACreate(r.Context(), request)
	response.CACertificate = response.Certificate
	rest.Response(w, response, err, http.StatusCreated, fmt.Sprintf("/v1/ca/%s", response.CAID))

}
