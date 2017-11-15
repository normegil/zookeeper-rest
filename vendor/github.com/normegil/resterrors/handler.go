package resterrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Handler struct {
	Definitions []ErrorDefinition
	DefaultCode int
}

func (h Handler) Handle(w http.ResponseWriter, e error) error {
	responseBody, err := h.ToResponse(e)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseBody.HTTPStatus)
	responseBodyJSON, err := json.Marshal(responseBody)
	if nil != err {
		return err
	}
	fmt.Fprintf(w, string(responseBodyJSON))
	return nil
}

func (h Handler) ToResponse(e error) (*ErrorResponse, error) {
	eWithCode := getErrWithCode(e, h.DefaultCode)

	for _, definition := range h.Definitions {
		if eWithCode.Code() == definition.Code {
			return definition.ToResponse(e)
		}
	}

	return nil, errors.Wrapf(e, "Could not find default error definition: the handler need to have both the default definition and default code")
}
