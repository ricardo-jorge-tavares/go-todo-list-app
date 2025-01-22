package models

import "errors"

var (
	ErrorBadRequest      = errors.New("bad request")
	ErrorInternalFailure = errors.New("internal failure")
	ErrorNotFound        = errors.New("not found")
)

type ErrorModel struct {
	errorType error
	code      string
	message   string
}

func NewError(err error, code, description string) error {
	return ErrorModel{
		errorType: err,
		code:      code,
		message:   description,
	}
}

func (e ErrorModel) Error() string {
	return e.message
}

func (e ErrorModel) GetType() error {
	return e.errorType
}
func (e ErrorModel) GetCode() string {
	return e.code
}
func (e ErrorModel) GetMessage() string {
	return e.message
}
func (e ErrorModel) HTTPStatus() string {
	return e.code
}
