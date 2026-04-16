package testhelpers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"do-your-dailies/server/internal/db"

	"gorm.io/gorm"
)

var schemaCounter uint64

func NewTransactionDB(t testing.TB, migrate func(database *gorm.DB) error) *gorm.DB {
	t.Helper()
	testDatabaseName := inferredDatabaseName()
	testSchemaName := schemaNameForTest(t.Name())

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
		fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=%s sslmode=disable", testDatabaseName),
		fmt.Sprintf("host=host.docker.internal port=5432 user=postgres password=postgres dbname=%s sslmode=disable", testDatabaseName),
	)

	var database *gorm.DB
	var lastErr error
	for _, dsn := range dsnCandidates {
		candidate, err := db.New(dsn)
		if err == nil {
			database = candidate
			break
		}
		lastErr = err
	}
	if database == nil {
		t.Fatalf("open postgres db: %v", lastErr)
	}

	if err := database.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", testSchemaName)).Error; err != nil {
		t.Fatalf("create postgres schema: %v", err)
	}

	schemaDSNCandidates := make([]string, 0, len(dsnCandidates))
	for _, dsn := range dsnCandidates {
		schemaDSNCandidates = append(schemaDSNCandidates, dsn+" search_path="+testSchemaName)
	}

	var schemaDatabase *gorm.DB
	for _, dsn := range schemaDSNCandidates {
		candidate, err := db.New(dsn)
		if err == nil {
			schemaDatabase = candidate
			break
		}
		lastErr = err
	}
	if schemaDatabase == nil {
		t.Fatalf("open postgres schema db: %v", lastErr)
	}

	baseSQLDB, err := database.DB()
	if err != nil {
		t.Fatalf("get postgres base sql db: %v", err)
	}
	if err := baseSQLDB.Close(); err != nil && !strings.Contains(err.Error(), "closed") {
		t.Fatalf("close postgres base db: %v", err)
	}

	if err := migrate(schemaDatabase); err != nil {
		t.Fatalf("migrate postgres db: %v", err)
	}

	tx := schemaDatabase.Begin()
	if tx.Error != nil {
		t.Fatalf("begin postgres transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil && err != gorm.ErrInvalidTransaction {
			t.Fatalf("rollback postgres transaction: %v", err)
		}
		sqlDB, err := schemaDatabase.DB()
		if err != nil {
			t.Fatalf("get postgres sql db: %v", err)
		}
		if err := sqlDB.Close(); err != nil && !strings.Contains(err.Error(), "closed") {
			t.Fatalf("close postgres db: %v", err)
		}
	})

	return tx
}

func inferredDatabaseName() string {
	_, filePath, _, ok := runtime.Caller(2)
	if !ok {
		return "dailies_test"
	}

	return databaseNameFromFile(filePath)
}

func databaseNameFromFile(filePath string) string {
	normalizedPath := strings.ReplaceAll(filePath, "\\", "/")
	normalized := strings.ToLower(filepath.ToSlash(normalizedPath))
	marker := "/internal/"
	index := strings.Index(normalized, marker)
	if index < 0 {
		return "dailies_test"
	}

	packagePath := filepath.ToSlash(filepath.Dir(normalized[index+1:]))
	if packagePath == "." || packagePath == "/" || packagePath == "" {
		return "dailies_test"
	}

	var builder strings.Builder
	for _, character := range packagePath {
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') {
			builder.WriteRune(character)
			continue
		}
		builder.WriteByte('_')
	}

	sanitized := strings.Trim(builder.String(), "_")
	if sanitized == "" {
		return "dailies_test"
	}

	fullName := "dailies_test_" + sanitized
	if len(fullName) <= 63 {
		return fullName
	}

	hash := sha1.Sum([]byte(fullName))
	suffix := hex.EncodeToString(hash[:])[:8]
	trimmed := fullName[:63-len(suffix)-1]
	trimmed = strings.TrimRight(trimmed, "_")
	if trimmed == "" {
		return "dailies_test_" + suffix
	}

	return trimmed + "_" + suffix
}

func ensureDatabaseExists(t testing.TB, adminDSNs []string, databaseName string) {
	t.Helper()

	var adminDB *gorm.DB
	var lastErr error
	for _, dsn := range adminDSNs {
		candidate, err := db.New(dsn)
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
			if !isDatabaseAlreadyCreatedError(err) {
				t.Fatalf("create postgres test database: %v", err)
			}
		}
	}
}

func isDatabaseAlreadyCreatedError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "already exists") || strings.Contains(message, "pg_database_datname_index")
}

func schemaNameForTest(testName string) string {
	base := strings.ToLower(testName)
	var builder strings.Builder
	for _, character := range base {
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') {
			builder.WriteRune(character)
			continue
		}
		builder.WriteByte('_')
	}

	cleaned := strings.Trim(builder.String(), "_")
	if cleaned == "" {
		cleaned = "test"
	}

	index := atomic.AddUint64(&schemaCounter, 1)
	name := "t_" + cleaned + "_" + strconv.FormatUint(index, 10)
	if len(name) <= 63 {
		return name
	}

	hash := sha1.Sum([]byte(name))
	suffix := hex.EncodeToString(hash[:])[:8]
	trimmed := name[:63-len(suffix)-1]
	trimmed = strings.TrimRight(trimmed, "_")
	if trimmed == "" {
		return "t_" + suffix
	}

	return trimmed + "_" + suffix
}
