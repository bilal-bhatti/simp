package skit

import (
	"io"
	"log"
	"testing"

	"github.com/pkg/errors"
)

func TestErrors(t *testing.T) {
	e := start()
	log.Printf("status +v: %+v\n", e)
	log.Printf("status #v: %#v\n", e)
	log.Printf("status  v: %v\n", e)
	log.Printf("status  s: %s\n", e)
	log.Printf("status  q: %q\n", e)

	e2 := errors.Wrap(e, "wrap error")
	log.Printf("wrap +v: %+v\n", e2)

	e3 := errors.Wrap(e2, "wrap error again")
	log.Printf("again +v: %+v\n", e3)

	log.Printf("cause: %+v", errors.Cause(e3))

	if ok, c, b := Status(e2); ok {
		log.Printf("status: %d - %v\n", c, b)
	}
}

func start() error {
	return do()
}

func do() error {
	return WithStatus(oops(), 500, []string{"one", "tow"})
}

func oops() error {
	return io.EOF
}
