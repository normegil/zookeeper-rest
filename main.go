package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

const PORT int = 8080

func init() {

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	})

	logrus.WithField("port", PORT).Info("Launch server")
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
