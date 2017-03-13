package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	Log *logrus.Entry
}

func (r *Router) Serve(port int) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Test")
	})

	r.Log.WithField("port", port).Info("Launching server")
	http.ListenAndServe(":"+strconv.Itoa(port), router)
}
