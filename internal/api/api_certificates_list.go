package api

import (
	"net/http"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// getCertificates GET /v1/ca/:caid/certificates
func (a *API) getCertificates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		response map[string]client.Certificate
		caID     string = ps.ByName("caid")
		parse    bool
		err      error
	)

	if r.URL.Query().Get("parse") == "true" {
		parse = true
	}

	response, err = a.srv.CertificateList(r.Context(), caID, parse)
	if len(response) == 0 {
		rest.NotFound(w, r)
		return
	}
	rest.Response(w, response, err, http.StatusOK, "")

}
