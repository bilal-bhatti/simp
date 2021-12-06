package skit

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

var (
	_ holder = &status{}
)

type tracer interface {
	error
	StackTrace() errors.StackTrace
}

type holder interface {
	error
	status() (int, interface{})
}

type status struct {
	err  error
	code int
	body interface{}
}

func (e *status) Error() string {
	return e.err.Error()
}

func (e *status) status() (int, interface{}) {
	return e.code, e.body
}

func (e *status) StackTrace() errors.StackTrace {
	if err, ok := e.err.(tracer); ok {
		return err.StackTrace()
	}

	return []errors.Frame{}
}

func (e *status) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, e.Error())
		e.StackTrace()[1:].Format(s, verb)
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *status) Unwrap() error {
	return e.err
}

func WithStatus(err error, code int, body interface{}) error {
	return &status{
		err:  errors.Errorf("%s - %+v", http.StatusText(code), err),
		code: code,
		body: body,
	}
}

func Status(err error) (bool, int, interface{}) {
	var sh holder
	if errors.As(err, &sh) {
		code, body := sh.status()
		return true, code, body
	}
	return false, 0, nil
}
