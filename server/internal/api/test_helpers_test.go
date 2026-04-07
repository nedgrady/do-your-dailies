package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"do-your-dailies/server/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPostgresTestApp(t *testing.T) (*Application, *gorm.DB) {
	t.Helper()
	tx := newCreateChoreTestTransactionDB(t)
	return New(tx), tx
}

func seedChore(t *testing.T, database *gorm.DB, name string, cadenceInDays int) models.Chore {
	t.Helper()

	chore := models.Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed chore: %v", err)
	}

	return chore
}

func breakChoreTable(t *testing.T, database *gorm.DB) {
	t.Helper()

	if err := database.Exec("DROP TABLE chores CASCADE").Error; err != nil {
		t.Fatalf("drop chores table: %v", err)
	}
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
		if err := sqlDB.Close(); err != nil && !strings.Contains(err.Error(), "closed") {
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

	var existing int64
	if lookupErr := adminDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", databaseName).Scan(&existing).Error; lookupErr != nil {
		t.Fatalf("ensure postgres test database: %v", lookupErr)
	}

	if existing == 0 {
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName)).Error; err != nil {
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
