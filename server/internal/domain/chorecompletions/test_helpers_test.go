package chorecompletions

import (
	"strconv"
	"testing"

	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func newTestRouter(choreStore chores.Store, store Store) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api/chore-completions", NewHandler(choreStore, store).RegisterRoutes)
	return router
}

func newPostgresTestRouter(t *testing.T) (*chi.Mux, *gorm.DB) {
	t.Helper()
	database := testhelpers.NewTransactionDB(t, migrate)
	return newTestRouter(chores.NewGormStore(database), NewGormStore(database)), database
}

func seedChore(t *testing.T, database *gorm.DB, name string, cadenceInDays int) chores.Chore {
	t.Helper()

	chore := chores.Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed chore: %v", err)
	}

	return chore
}

func itoa(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(&chores.Chore{}, &ChoreCompletion{})
}
