package testhelpers

import (
	"bytes"
	"net/http"
)

type RecordingResponseWriter struct {
	header                 http.Header
	StatusCode             int
	HeaderWritten          bool
	WriteBeforeWriteHeader bool
	Body                   bytes.Buffer
}

func NewRecordingResponseWriter() *RecordingResponseWriter {
	return &RecordingResponseWriter{header: make(http.Header)}
}

func (writer *RecordingResponseWriter) Header() http.Header {
	return writer.header
}

func (writer *RecordingResponseWriter) WriteHeader(statusCode int) {
	writer.HeaderWritten = true
	writer.StatusCode = statusCode
}

func (writer *RecordingResponseWriter) Write(payload []byte) (int, error) {
	if !writer.HeaderWritten {
		writer.WriteBeforeWriteHeader = true
	}

	return writer.Body.Write(payload)
}
