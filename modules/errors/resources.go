package errors

import (
	"bytes"
	"encoding/csv"
	"net/url"
	"strconv"

	urlFormat "github.com/normegil/formats/url"
	"github.com/pkg/errors"
)

const baseAddress = "http://example.com/"

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
var predefinedErrors []ErrorResponse

func init() {
	predefinedErrors = loadErrorRessources("assets/errors.csv")
}

func loadErrorRessources(path string) []ErrorResponse {
	assetBytes, err := Asset(path)
	if nil != err {
		panic(errors.Wrapf(err, "Loading assets %s", path))
	}
	content, err := csv.NewReader(bytes.NewReader(assetBytes)).ReadAll()
	if nil != err {
		panic(errors.Wrapf(err, "Parsing content of %s", path))
	}
	loadedResponses := make([]ErrorResponse, 0)
	for _, row := range content {
		code, err := strconv.Atoi(row[0])
		if nil != err {
			panic(errors.Wrapf(err, "Converting code %s", row[0]))
		}
		httpStatus, err := strconv.Atoi(row[1])
		if nil != err {
			panic(errors.Wrapf(err, "Converting HTTP Status code %s", row[1]))
		}
		urlStr := baseAddress + strconv.Itoa(code)
		moreInfo, err := url.Parse(urlStr)
		if nil != err {
			panic(errors.Wrapf(err, "Parsing MoreInfo URL %s", urlStr))
		}
		loadedResponses = append(loadedResponses, ErrorResponse{
			Code:       code,
			HTTPStatus: httpStatus,
			MoreInfo:   urlFormat.URL{moreInfo},
			Message:    row[2],
		})
	}
	return loadedResponses
}
