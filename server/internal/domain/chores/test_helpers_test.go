package chores

import (
	"testing"

	"do-your-dailies/server/internal/testhelpers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func newTestRouter(store Store) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api/chores", NewHandler(store).RegisterRoutes)
	return router
}

func newPostgresTestRouter(t *testing.T) (*chi.Mux, *gorm.DB) {
	t.Helper()
	database := testhelpers.NewTransactionDB(t, migrate)
	return newTestRouter(NewGormStore(database)), database
}

func seedChore(t *testing.T, database *gorm.DB, name string, cadenceInDays int) Chore {
	t.Helper()

	chore := Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed chore: %v", err)
	}

	return chore
}

func breakChoreTable(t *testing.T, database *gorm.DB) {
	t.Helper()
	testhelpers.BreakTable(t, database, "chores")
}

func newCreateChoreTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testhelpers.NewTransactionDB(t, migrate)
}

func newCreateChoreTestTransactionDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testhelpers.NewTransactionDB(t, migrate)
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(&Chore{})
}
