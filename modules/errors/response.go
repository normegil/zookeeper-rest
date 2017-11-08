package errors

import (
	"fmt"
	"net/url"
	"time"

	errorFormat "github.com/normegil/formats/error"
	timeFormat "github.com/normegil/formats/time"
	urlFormat "github.com/normegil/formats/url"
)

const DEFAULT_CODE = 500

func toResponse(e error) *ErrorResponse {
	eWithCode := getErrWithCode(e)
	fmt.Printf("%+v\n", eWithCode)
	eMarshallable, isMarshable := e.(marshableError)
	if !isMarshable {
		eMarshallable = errorFormat.Error{e.Error()}
	}

	for _, defResp := range predefinedErrors {
		if eWithCode.Code() == defResp.Code {
			return &ErrorResponse{
				Code:       defResp.Code,
				HTTPStatus: defResp.HTTPStatus,
				Message:    defResp.Message,
				MoreInfo:   defResp.MoreInfo,
				Time:       timeFormat.Time(time.Now()),
				Err:        eMarshallable,
			}
		}
	}

	moreInfo, err := url.Parse("http://example.com/5000")
	if nil != err {
		panic(err)
	}
	return &ErrorResponse{
		Code:       50000,
		HTTPStatus: 500,
		Err:        errorFormat.Error{e.Error()},
		MoreInfo:   urlFormat.URL{moreInfo},
		Time:       timeFormat.Time(time.Now()),
		Message:    "Error was not found in the error ressources. Generated a default error.",
	}
}
