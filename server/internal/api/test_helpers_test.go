package api

import (
	"bytes"
	"net/http"

	"do-your-dailies/server/internal/chores"
)

type mockChoreStore struct {
	listFn   func() ([]chores.Chore, error)
	createFn func(chores.CreateRequest) (chores.Chore, error)
	getFn    func(uint) (chores.Chore, error)
	updateFn func(uint, chores.UpdateRequest) (chores.Chore, error)
	deleteFn func(uint) error
}

func (m *mockChoreStore) List() ([]chores.Chore, error) { return m.listFn() }
func (m *mockChoreStore) Create(req chores.CreateRequest) (chores.Chore, error) {
	return m.createFn(req)
}
func (m *mockChoreStore) Get(id uint) (chores.Chore, error) { return m.getFn(id) }
func (m *mockChoreStore) Update(id uint, req chores.UpdateRequest) (chores.Chore, error) {
	return m.updateFn(id, req)
}
func (m *mockChoreStore) Delete(id uint) error { return m.deleteFn(id) }

func newAppWithStore(mock chores.Store) *Application {
	app := New(nil)
	app.ChoreStore = mock
	app.Router = app.setupRoutes()
	return app
}

type recordingResponseWriter struct {
	header                 http.Header
	statusCode             int
	headerWritten          bool
	writeBeforeWriteHeader bool
	body                   bytes.Buffer
}

func newRecordingResponseWriter() *recordingResponseWriter {
	return &recordingResponseWriter{header: make(http.Header)}
}

func (w *recordingResponseWriter) Header() http.Header {
	return w.header
}

func (w *recordingResponseWriter) WriteHeader(statusCode int) {
	w.headerWritten = true
	w.statusCode = statusCode
}

func (w *recordingResponseWriter) Write(payload []byte) (int, error) {
	if !w.headerWritten {
		w.writeBeforeWriteHeader = true
	}

	return w.body.Write(payload)
}
