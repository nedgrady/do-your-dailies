package migrations

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"do-your-dailies/server/internal/db"

	"gorm.io/gorm"
)

func TestMigrateCreatesChoresTable(t *testing.T) {
	t.Parallel()
	schemaName := testSchemaName(t)
	adminDB := newTestDatabase(t, []string{
		"host=localhost port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
		"host=host.docker.internal port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
	})
	t.Cleanup(func() {
		closeSQLDB(t, adminDB)
	})
	createSchema(t, adminDB, schemaName)
	t.Cleanup(func() {
		dropSchema(t, adminDB, schemaName)
	})

	schemaDB := newTestDatabase(t, testSchemaDSNCandidates(schemaName))
	t.Cleanup(func() {
		closeSQLDB(t, schemaDB)
	})

	if err := Migrate(schemaDB); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	var count int64
	if err := adminDB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'chores'", schemaName).Scan(&count).Error; err != nil {
		t.Fatalf("count chores table: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected chores table to exist, got %d", count)
	}
}

func TestMigrateCreatesChoreCompletionsTable(t *testing.T) {
	t.Parallel()
	schemaName := testSchemaName(t)
	adminDB := newTestDatabase(t, []string{
		"host=localhost port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
		"host=host.docker.internal port=5432 user=postgres password=postgres dbname=dailies_test sslmode=disable",
	})
	t.Cleanup(func() {
		closeSQLDB(t, adminDB)
	})
	createSchema(t, adminDB, schemaName)
	t.Cleanup(func() {
		dropSchema(t, adminDB, schemaName)
	})

	schemaDB := newTestDatabase(t, testSchemaDSNCandidates(schemaName))
	t.Cleanup(func() {
		closeSQLDB(t, schemaDB)
	})

	if err := Migrate(schemaDB); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	var count int64
	if err := adminDB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'chore_completions'", schemaName).Scan(&count).Error; err != nil {
		t.Fatalf("count chore_completions table: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected chore_completions table to exist, got %d", count)
	}
}

func newTestDatabase(t *testing.T, dsnCandidates []string) *gorm.DB {
	t.Helper()
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
		t.Fatalf("open test database: %v", lastErr)
	}

	return database
}

func testSchemaDSNCandidates(schemaName string) []string {
	candidates := []string{}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		candidates = append(candidates, dsn)
	}

	candidates = append(candidates,
		fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=dailies_test search_path=%s sslmode=disable", schemaName),
		fmt.Sprintf("host=host.docker.internal port=5432 user=postgres password=postgres dbname=dailies_test search_path=%s sslmode=disable", schemaName),
	)

	return candidates
}

func createSchema(t *testing.T, database *gorm.DB, schemaName string) {
	t.Helper()
	if err := database.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName)).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
}

func dropSchema(t *testing.T, database *gorm.DB, schemaName string) {
	t.Helper()
	if err := database.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", schemaName)).Error; err != nil {
		t.Fatalf("drop schema: %v", err)
	}
}

func closeSQLDB(t *testing.T, database *gorm.DB) {
	t.Helper()
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	if err := sqlDB.Close(); err != nil && !strings.Contains(err.Error(), "closed") {
		t.Fatalf("close sql db: %v", err)
	}
}

func testSchemaName(t *testing.T) string {
	t.Helper()
	replacer := strings.NewReplacer("/", "_", "-", "_")
	return fmt.Sprintf("test_%s_%d", replacer.Replace(strings.ToLower(t.Name())), time.Now().UnixNano())
}
