package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"do-your-dailies/server/internal/models"
	"do-your-dailies/server/internal/store"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockChoreStore struct {
	listFn   func() ([]models.Chore, error)
	createFn func(models.CreateChoreRequest) (models.Chore, error)
	getFn    func(uint) (models.Chore, error)
	updateFn func(uint, models.UpdateChoreRequest) (models.Chore, error)
	deleteFn func(uint) error
}

func (m *mockChoreStore) List() ([]models.Chore, error) { return m.listFn() }
func (m *mockChoreStore) Create(req models.CreateChoreRequest) (models.Chore, error) {
	return m.createFn(req)
}
func (m *mockChoreStore) Get(id uint) (models.Chore, error) { return m.getFn(id) }
func (m *mockChoreStore) Update(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
	return m.updateFn(id, req)
}
func (m *mockChoreStore) Delete(id uint) error { return m.deleteFn(id) }

func newAppWithStore(mock store.ChoreStore) *Application {
	app := New(nil)
	app.ChoreStore = mock
	app.Router = app.setupRoutes()
	return app
}

func newCreateChoreTestDB(t *testing.T) *gorm.DB {
	return newCreateChoreTestTransactionDB(t)
}

func newCreateChoreTestTransactionDB(t *testing.T) *gorm.DB {
	t.Helper()
	testDatabaseName := "dailies_test"

	adminDSNCandidates := []string{
		"host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
		"host=host.docker.internal port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
	}
	ensureDatabaseExists(t, adminDSNCandidates, testDatabaseName)

	dsnCandidates := []string{}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		dsnCandidates = append(dsnCandidates, dsn)
	}
	dsnCandidates = append(dsnCandidates,
		"host=localhost port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
		"host=host.docker.internal port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
	)

	var database *gorm.DB
	var lastErr error
	for _, dsn := range dsnCandidates {
		candidate, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			database = candidate
			break
		}
		lastErr = err
	}
	if database == nil {
		t.Fatalf("open postgres db: %v", lastErr)
	}

	if err := database.AutoMigrate(&models.Chore{}, &models.ChoreCompletion{}); err != nil {
		t.Fatalf("migrate postgres db: %v", err)
	}

	tx := database.Begin()
	if tx.Error != nil {
		t.Fatalf("begin postgres transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil && err != gorm.ErrInvalidTransaction {
			t.Fatalf("rollback postgres transaction: %v", err)
		}
		sqlDB, err := database.DB()
		if err != nil {
			t.Fatalf("get postgres sql db: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("close postgres db: %v", err)
		}
	})

	return tx
}

func ensureDatabaseExists(t *testing.T, adminDSNs []string, databaseName string) {
	t.Helper()

	var adminDB *gorm.DB
	var lastErr error
	for _, dsn := range adminDSNs {
		candidate, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			adminDB = candidate
			break
		}
		lastErr = err
	}
	if adminDB == nil {
		t.Fatalf("open postgres admin db: %v", lastErr)
	}

	adminSQLDB, err := adminDB.DB()
	if err != nil {
		t.Fatalf("get postgres admin sql db: %v", err)
	}
	defer func() {
		if closeErr := adminSQLDB.Close(); closeErr != nil {
			t.Fatalf("close postgres admin db: %v", closeErr)
		}
	}()

	if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName)).Error; err != nil {
		var existing int64
		if lookupErr := adminDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", databaseName).Scan(&existing).Error; lookupErr != nil {
			t.Fatalf("ensure postgres test database: %v", lookupErr)
		}
		if existing == 0 {
			t.Fatalf("create postgres test database: %v", err)
		}
	}
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
