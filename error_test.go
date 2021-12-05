package skit

import (
	"errors"
	"io"
	"log"
	"testing"
)

func TestErrors(t *testing.T) {
	log.Print(io.EOF)
	log.Printf("%v\n", New("an error"))
	log.Printf("%+v\n", Newf("an error with string: %s", "hi"))
	e := Wrap(io.EOF, "wrapped eof error")
	log.Printf("%+v\n", e)
	log.Printf("%v\n", e)
	log.Printf("%s:%v\n", e, e)

	log.Println("---")
	e = Wrap(io.EOF, "unexpected eof")
	e = WithStatus(e, 404, []string{"one", "tow"})

	if r, ok := e.(StatusHolder); ok {
		code, body := r.Status()
		log.Printf("body 1: %d - %v\n", code, body)
	}

	e = Wrap(e, "nest the response")
	log.Printf("+v: %+v\n", e)
	log.Printf(" v: %v\n", e)

	if r, ok := e.(StatusHolder); ok {
		code, body := r.Status()
		log.Printf("body 2: %d - %v\n", code, body)
	}

	var sh StatusHolder

	if errors.As(e, &sh) {
		c, b := sh.Status()
		log.Printf("body 3: %d - %v\n", c, b)
		log.Printf("body 3: %+v\n", sh)
	}
}
