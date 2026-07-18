package chorecompletions

import (
	"strconv"
	"testing"

	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	chorequeue "do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func newTestRouter(choreStore chores.Store, store Store, choreQueueStore chorequeue.Store, choreInQueueCompletionStore choreinqueuecompletion.Store, db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api/chore-completions", NewHandler(choreStore, store, choreQueueStore, choreInQueueCompletionStore, db).RegisterRoutes)
	return router
}

func newPostgresTestRouter(t *testing.T) (*chi.Mux, *gorm.DB) {
	t.Helper()
	database := testhelpers.NewTransactionDB(t, migrate)
	return newTestRouter(chores.NewGormStore(database), NewGormStore(database), chorequeue.NewGormStore(database), choreinqueuecompletion.NewGormStore(database), database), database
}

func seedChore(t *testing.T, database *gorm.DB, name string, cadenceInDays int) models.Chore {
	t.Helper()

	chore := models.Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed chore: %v", err)
	}

	return chore
}

func itoa(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(&models.Chore{}, &models.ChoreCompletion{})
}
