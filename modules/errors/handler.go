package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log logrus.FieldLogger
}

func (h Handler) Handle(w http.ResponseWriter, e error) {
	log := h.Log
	stacks := stacks(e)
	if len(stacks) > 0 {
		log = log.WithField("errorStack", stacks[0])
	}
	log.WithError(e).Error("Error while processing request")

	responseBody := toResponse(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseBody.HTTPStatus)
	responseBodyJSON, err := json.Marshal(responseBody)
	if nil != err {
		fmt.Fprint(w, "An Error ("+err.Error()+") happened when trying to marshall Error to JSON. "+responseBody.String())
		h.Log.WithError(err).Error("An error happened while trying to marshall an other error")
		return
	}
	fmt.Fprint(w, string(responseBodyJSON))
	log.WithField("headers", w.Header()).Debug("Headers of error response")
}
