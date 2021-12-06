package skit

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func Success(w http.ResponseWriter, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logError("", errors.Errorf("failed to encode data: %+v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		logError("", errors.Errorf("failed to write http response: %+v", err))
	}
}

type errorResponse struct {
	Code    int         `json:"code"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func Failure(w http.ResponseWriter, err error) {
	logError("failure response", err)

	resp := &errorResponse{
		Code:  http.StatusInternalServerError,
		Error: err.Error(),
	}

	if ok, c, b := Status(err); ok {
		resp.Code = c
		resp.Details = b
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(resp); err != nil {
		http.Error(w, http.StatusText(resp.Code), resp.Code)
		logError("", errors.Errorf("failed to encode data: %+v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.Code)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		logError("", errors.Errorf("failed to write http response: %+v", err))
		return
	}
}

func logError(msg string, err error) {
	if msg != "" {
		log.Println(msg)
	}

	if err == nil {
		return
	}

	log.Printf("stacktrace: %+v", err)
}
