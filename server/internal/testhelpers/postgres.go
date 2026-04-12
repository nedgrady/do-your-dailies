package testhelpers

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"do-your-dailies/server/internal/db"

	"gorm.io/gorm"
)

func NewTransactionDB(t testing.TB, migrate func(database *gorm.DB) error) *gorm.DB {
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

	if err := migrate(database); err != nil {
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

func BreakTable(t testing.TB, database *gorm.DB, tableName string) {
	t.Helper()

	if err := database.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", tableName)).Error; err != nil {
		t.Fatalf("drop %s table: %v", tableName, err)
	}
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
			t.Fatalf("create postgres test database: %v", err)
		}
	}
}
