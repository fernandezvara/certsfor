package api

import (
	"net/http"
	"strconv"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// getCertificate GET /v1/ca/:caid/certificates/:cn
func (a *API) getCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		response        client.Certificate
		caID            string = ps.ByName("caid")
		cn              string = ps.ByName("cn") // certificate common name
		remainingString string
		remaining       int
		parse           bool
		err             error
	)

	remainingString = r.URL.Query().Get("renew")
	if remainingString != "" {
		remaining, err = strconv.Atoi(remainingString)
		if err != nil {
			rest.BadRequest(w, r, "renew value not allowed")
			return
		}
	}

	if r.URL.Query().Get("parse") == "true" {
		parse = true
	}

	response, err = a.srv.CertificateGet(r.Context(), caID, cn, remaining, parse)
	rest.Response(w, response, err, http.StatusOK, "")

}
