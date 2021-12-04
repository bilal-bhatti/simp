package simp

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	_ error = &AppError{}
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type ResponseHolder interface {
	StatusCode() int
	Body() interface{}
}

type AppError struct {
	err  error
	code int
	body interface{}
}

func New(msg string) *AppError {
	return &AppError{
		err:  errors.New(msg),
		code: http.StatusInternalServerError,
	}
}

func Wrap(err error, msg string) *AppError {
	return &AppError{
		err:  errors.Wrap(err, msg),
		code: http.StatusInternalServerError,
	}
}

func (e *AppError) StackTrace() errors.StackTrace {
	if err, ok := e.err.(StackTracer); ok {
		return err.StackTrace()
	}

	// use runtime to build the frames, maybe?
	return []errors.Frame{}
}

func (e *AppError) StatusCode() int {
	return e.code
}

func (e *AppError) Body() interface{} {
	return e.body
}

func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) WithStatusCode(code int) *AppError {
	e.code = code
	return e
}

func (e *AppError) WithBody(body interface{}) *AppError {
	e.body = body
	return e
}
