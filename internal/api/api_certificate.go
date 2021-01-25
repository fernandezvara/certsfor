package api

import (
	"net/http"

	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// putCertificate PUT /v1/ca/:caid/certificates/%s
func (a *API) putCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		request  client.APICertificateRequest
		response client.Certificate
		caID     string = ps.ByName("caid")
		cn       string = ps.ByName("cn") // certificate common name
		err      error
	)

	err = rest.GetFromBody(r, &request)
	if err != nil {
		rest.BadRequest(w, r, "")
		return
	}

	if cn == "ca" {
		rest.ErrorResponse(w, http.StatusConflict, "CA Certificate cannot be overwritten")
		return
	}

	if request.DN.CN == "" {
		request.DN.CN = cn
	}

	if request.DN.CN != cn {
		rest.ErrorResponse(w, http.StatusConflict, "Common Name on certificate does not match")
		return
	}

	response.Request = request

	response.CACertificate, response.Certificate, response.Key, err = a.srv.CertificateSet(r.Context(), caID, request)
	rest.Response(w, response, err, http.StatusOK, "")

}
