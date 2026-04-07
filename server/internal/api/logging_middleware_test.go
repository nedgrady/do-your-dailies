package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

var loggerMutationLock sync.Mutex

func TestRouterLogsHealthRequests(t *testing.T) {
	t.Parallel()

	loggerMutationLock.Lock()
	t.Cleanup(func() {
		loggerMutationLock.Unlock()
	})

	originalLogger := middleware.DefaultLogger
	var logOutput bytes.Buffer
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(&logOutput, "", 0),
		NoColor: true,
	})
	t.Cleanup(func() {
		middleware.DefaultLogger = originalLogger
	})

	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.NotEmpty(t, logOutput.String())
}
