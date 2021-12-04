package simp

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

func Success(w http.ResponseWriter, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logError(0, "failed to encode JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		logError(0, "failed to write http response", err)
	}
}

type errorResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

func Error(w http.ResponseWriter, err error) {
	defer func() {
		logError(2, "", err)
	}()

	resp := &errorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}

	if holder, ok := err.(ResponseHolder); ok {
		resp.Code = holder.StatusCode()
		if holder.Body() != nil {
			resp.Body = holder.Body()
		}
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(resp); err != nil {
		http.Error(w, http.StatusText(resp.Code), resp.Code)
		logError(0, "failed to encode error response to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.Code)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		logError(0, "failed to write error response", err)
		return
	}
}

func logError(skip int, msg string, err error) {
	if msg != "" {
		log.Println(msg)
	}

	if err == nil {
		return
	}

	// skip callers
	pc, _, line, _ := runtime.Caller(skip)

	log.Printf("calling function: %s\n", runtime.FuncForPC(pc).Name())
	log.Printf("line: %d\n", line)
	log.Println("stacktrace:")

	if err, ok := err.(StackTracer); ok {
		st := err.StackTrace()[skip:]
		for _, frame := range st {
			log.Printf("%+s:%d", frame, frame)
		}
	}
}
