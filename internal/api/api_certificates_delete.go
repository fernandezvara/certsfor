package api

import (
	"net/http"

	"github.com/fernandezvara/rest"
	"github.com/julienschmidt/httprouter"
)

// deleteCertificate GET /v1/ca/:caid/certificates/:cn
func (a *API) deleteCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		caID string = ps.ByName("caid")
		cn   string = ps.ByName("cn") // certificate common name
		err  error
	)

	_, err = a.srv.CertificateDelete(r.Context(), caID, cn)
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
