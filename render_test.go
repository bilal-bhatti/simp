package skit

import (
	"log"
	"net/http/httptest"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		in   interface{}
		code int
	}{
		{
			in:   []string{"one", "two"},
			code: 200,
		},
		{
			in:   make(chan int),
			code: 500,
		},
	}

	for _, test := range tests {
		resp := httptest.NewRecorder()

		Success(resp, test.in)

		log.Print(string(resp.Body.Bytes()))
		log.Println(resp.Header())
		log.Println(resp.Result().Status)
		log.Println(resp.Result().StatusCode)

		if test.code != resp.Result().StatusCode {
			t.Errorf("code mismatch - expected: %d, actual: %d", test.code, resp.Result().StatusCode)
		}
	}
}
