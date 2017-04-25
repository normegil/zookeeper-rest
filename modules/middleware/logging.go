package middleware

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

func RequestLogger(log *logrus.Entry, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("method", r.Method).WithField("URL", r.Context().Value(FULL_URL_KEY)).Info("Request received")
		log.WithField("Request", r).Debug("Request received (Details)")
		h.ServeHTTP(w, r)
	})
}
