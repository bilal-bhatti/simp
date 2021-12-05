package skit

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	_ error        = &wrap{}
	_ StackTracer  = &wrap{}
	_ error        = &status{}
	_ StatusHolder = &status{}
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type StatusHolder interface {
	Status() (int, interface{})
}

type wrap struct {
	err error
}

func New(msg string) error {
	return &wrap{
		err: errors.New(msg),
	}
}

func Newf(msg string, args ...interface{}) error {
	return &wrap{
		err: errors.Errorf(msg, args...),
	}
}

func Wrap(err error, msg string) error {
	return &wrap{
		err: errors.Wrap(err, msg),
	}
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return &wrap{
		err: errors.Wrapf(err, msg, args...),
	}
}

func (e *wrap) StackTrace() errors.StackTrace {
	if err, ok := e.err.(StackTracer); ok {
		return err.StackTrace()
	}

	return []errors.Frame{}
}

func (e *wrap) Error() string {
	return e.err.Error()
}

func (e *wrap) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		e.StackTrace().Format(s, verb)
	case 's':
		io.WriteString(s, e.Error())
	}
}

func (e *wrap) Unwrap() error {
	return e.err
}

type status struct {
	err  error
	code int
	body interface{}
}

func (e *status) Error() string {
	return e.err.Error()
}

func (e *status) Status() (int, interface{}) {
	return e.code, e.body
}

func (e *status) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if t, ok := e.err.(StackTracer); ok {
			t.StackTrace().Format(s, verb)
		} else {
			if s.Flag('+') {
				fmt.Fprintf(s, "%+v", t)
			} else {
				fmt.Fprintf(s, "%v", t)
			}
		}
	case 's':
		io.WriteString(s, e.Error())
	}
}

func (e *status) Unwrap() error {
	return e.err
}

func WithStatus(err error, code int, body interface{}) error {
	return &status{
		err:  errors.WithStack(err),
		code: code,
		body: body,
	}
}
