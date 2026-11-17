package idempotency

import (
	"bytes"
	"net/http"
)

// recorder passes the handler's response through to the client while keeping a
// copy of the status, content type and body so they can be stored and replayed.
type recorder struct {
	http.ResponseWriter
	status      int
	body        bytes.Buffer
	wroteHeader bool
}

func newRecorder(w http.ResponseWriter) *recorder {
	return &recorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *recorder) WriteHeader(status int) {
	if !r.wroteHeader {
		r.status = status
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(p []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	r.body.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *recorder) contentType() string {
	return r.Header().Get("Content-Type")
}
