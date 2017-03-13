package middleware

import (
	"net/http"

	"github.com/normegil/zookeeper-rest/model"
)

func RequestLogger(env model.Env, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env.Log().WithField("URL", r.URL).Info("Request received")
		env.Log().WithField("Request", r).Debug("Request received (Details)")
		h.ServeHTTP(w, r)
	})
}
