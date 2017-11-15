package resterrors

import (
	"net/url"
	"strconv"
	"time"

	"encoding/json"

	errorFormat "github.com/normegil/formats/error"
	timeFormat "github.com/normegil/formats/time"
	urlFormat "github.com/normegil/formats/url"
	"github.com/pkg/errors"
)

type ErrWithCode interface {
	Code() int
	error
}

type ErrWithCodeImpl struct {
	error
	code int
}

func (e ErrWithCodeImpl) Code() int {
	return e.code
}

func NewErrWithCode(code int, e error) ErrWithCode {
	return &ErrWithCodeImpl{
		code:  code,
		error: e,
	}
}

type ErrorDefinition struct {
	Code       int
	HTTPStatus int
	MoreInfo   string
	Message    string
}

func (d ErrorDefinition) ToResponse(err error) (*ErrorResponse, error) {
	eMarshallable, isMarshable := err.(marshableError)
	if !isMarshable {
		eMarshallable = errorFormat.Error{err.Error()}
	}
	moreInfoURL, err := url.Parse(d.MoreInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "Parsing %s as URL", d.MoreInfo)
	}
	moreInfo := urlFormat.URL{moreInfoURL}
	return &ErrorResponse{
		HTTPStatus: d.HTTPStatus,
		Code:       d.Code,
		Message:    d.Message,
		MoreInfo:   moreInfo,
		Time:       timeFormat.Time(time.Now()),
		Err:        eMarshallable,
	}, nil
}

type ErrorResponse struct {
	HTTPStatus int             `json:"http status"`
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	MoreInfo   urlFormat.URL   `json:"more info"`
	Time       timeFormat.Time `json:"time"`
	Err        marshableError  `json:"error"`
}

func (e ErrorResponse) String() string {
	return "[Status HTTP:" + strconv.Itoa(e.HTTPStatus) + ";Code:" + strconv.Itoa(e.Code) + ";URL:" + e.MoreInfo.RawPath + ";Time:" + e.Time.String() + ";Msg:" + e.Message + ";Err:" + e.Err.Error() + "]"
}

func (e *ErrorResponse) UnmarshalJSON(b []byte) error {
	objRawMessages := make(map[string]*json.RawMessage)
	err := json.Unmarshal(b, &objRawMessages)
	if err != nil {
		return errors.Wrap(err, "Could not parse bytes into json RawMessages")
	}

	if err = json.Unmarshal([]byte(*objRawMessages["http status"]), &e.HTTPStatus); err != nil {
		return errors.Wrap(err, "Parsing HTTP Status")
	}
	if err = json.Unmarshal([]byte(*objRawMessages["code"]), &e.Code); err != nil {
		return errors.Wrap(err, "Parsing Code")
	}
	if err = json.Unmarshal([]byte(*objRawMessages["message"]), &e.Message); err != nil {
		return errors.Wrap(err, "Parsing Message")
	}
	if err = json.Unmarshal([]byte(*objRawMessages["more info"]), &e.MoreInfo); err != nil {
		return errors.Wrap(err, "Parsing MoreInfo")
	}
	if err = json.Unmarshal([]byte(*objRawMessages["time"]), &e.Time); err != nil {
		return errors.Wrap(err, "Parsing Time")
	}
	errForResponse := errorFormat.Error{}
	err = json.Unmarshal([]byte(*objRawMessages["error"]), &errForResponse)
	if err != nil {
		return errors.Wrap(err, "Parsing Error (response field)")
	}
	e.Err = errForResponse
	return nil
}

type marshableError interface {
	json.Marshaler
	error
}
