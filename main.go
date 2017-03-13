package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

const PORT int = 8080

var log *logrus.Entry = logrus.NewEntry(logrus.New())

func init() {
	log = log.WithField("Test", "123")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	})

	log.WithField("port", PORT).Info("Launch server")
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
