package errors

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
