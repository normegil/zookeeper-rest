package errors

import (
	"bytes"
	"encoding/csv"
	"net/url"
	"strconv"
	"time"

	"github.com/normegil/zookeeper-rest/modules/formats/json"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
var defaultResponses []response

func init() {
	assetBytes, err := Asset("assets/errors.csv")
	if nil != err {
		panic(err)
	}
	content, err := csv.NewReader(bytes.NewReader(assetBytes)).ReadAll()
	if nil != err {
		panic(err)
	}
	for _, row := range content {
		code, err := strconv.Atoi(row[0])
		if nil != err {
			panic(err)
		}
		httpStatus, err := strconv.Atoi(row[1])
		if nil != err {
			panic(err)
		}
		moreInfo, err := url.Parse(row[2])
		if nil != err {
			panic(err)
		}
		defaultResponses = append(defaultResponses, response{
			Code:             code,
			HTTPStatus:       httpStatus,
			MoreInfo:         json.JSONURL(*moreInfo),
			Message:          row[3],
			DeveloperMessage: row[4],
		})
	}
}

type response struct {
	HTTPStatus       int            `json:"http status"`
	Code             int            `json:"code"`
	Message          string         `json:"message"`
	DeveloperMessage string         `json:"developper message"`
	MoreInfo         json.JSONURL   `json:"more info"`
	Time             json.JSONTime  `json:"time"`
	Err              json.ErrorJSON `json:"error"`
}

const DEFAULT_CODE = 500

func newResponse(e error) *response {
	code := DEFAULT_CODE
	if eWithCode, ok := e.(ErrWithCode); ok {
		code = eWithCode.Code()
	}

	for _, defResp := range defaultResponses {
		if code == defResp.Code {
			return &response{
				Code:             defResp.Code,
				HTTPStatus:       defResp.HTTPStatus,
				Message:          defResp.Message,
				DeveloperMessage: defResp.DeveloperMessage,
				MoreInfo:         defResp.MoreInfo,
				Time:             json.JSONTime(time.Now()),
				Err:              json.ErrorJSON{e},
			}
		}
	}

	moreInfo, err := url.Parse("http://example.com/5000")
	if nil != err {
		panic(err)
	}

	return &response{
		Code:             50000,
		HTTPStatus:       500,
		Err:              json.ErrorJSON{e},
		MoreInfo:         json.JSONURL(*moreInfo),
		Time:             json.JSONTime(time.Now()),
		Message:          "An unrecognized error occured on the server",
		DeveloperMessage: "Error was not found in the error ressources. Generated a default error.",
	}
}

func (e response) String() string {
	return "[Status HTTP:" + strconv.Itoa(e.HTTPStatus) + ";Code:" + strconv.Itoa(e.Code) + ";URL:" + e.MoreInfo.RawPath + ";Time:" + e.Time.String() + ";Msg:" + e.Message + ";DevMsg:" + e.DeveloperMessage + ";Err:" + e.Err.Error() + "]"
}
