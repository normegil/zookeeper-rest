package middleware

import (
	"context"
	"net/http"
)

const FULL_URL_KEY = "RequestURL"

func URLContructor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), FULL_URL_KEY, "http://"+r.Host+r.URL.String())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
