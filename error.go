package skit

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	_ error = &appError{}
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type ResponseHolder interface {
	StatusCode() int
	Body() interface{}
}

type appError struct {
	err  error
	code int
	body interface{}
}

func New(msg string) ResponseHolder {
	return &appError{
		err:  errors.New(msg),
		code: http.StatusInternalServerError,
	}
}

func Wrap(err error, msg string) ResponseHolder {
	return &appError{
		err:  errors.Wrap(err, msg),
		code: http.StatusInternalServerError,
	}
}

func (e *appError) StackTrace() errors.StackTrace {
	if err, ok := e.err.(StackTracer); ok {
		return err.StackTrace()
	}

	// use runtime to build the frames, maybe?
	return []errors.Frame{}
}

func (e *appError) StatusCode() int {
	return e.code
}

func (e *appError) Body() interface{} {
	return e.body
}

func (e *appError) Error() string {
	return e.err.Error()
}

func (e *appError) WithStatusCode(code int) ResponseHolder {
	e.code = code
	return e
}

func (e *appError) WithBody(body interface{}) ResponseHolder {
	e.body = body
	return e
}
