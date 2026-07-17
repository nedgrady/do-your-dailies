package chorequeue

import (
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func newTestRouter(store Store, now func() time.Time, choreInQueueCompletionStore choreinqueuecompletion.Store) *chi.Mux {
	router := chi.NewRouter()
	handler := NewHandler(store, choreInQueueCompletionStore)
	handler.Now = now
	router.Route("/api/chore-queue", handler.RegisterRoutes)
	return router
}

func newPostgresTestRouter(t *testing.T) (*chi.Mux, *gorm.DB) {
	t.Helper()
	database := testhelpers.NewTransactionDB(t, migrate)
	if err := database.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ChoreCompletion{}).Error; err != nil {
		t.Fatalf("clear chore completions: %v", err)
	}
	if err := database.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Chore{}).Error; err != nil {
		t.Fatalf("clear chores: %v", err)
	}
	return newTestRouter(NewGormStore(database), fixedNow, choreinqueuecompletion.NewGormStore(database)), database
}

func seedChoreWithCompletion(t *testing.T, database *gorm.DB, name string, cadenceInDays int, createdDayOffset int, completedDayOffset int) {
	t.Helper()

	chore := models.Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed chore: %v", err)
	}
	createdAt := day(createdDayOffset)
	if err := database.Model(&chore).Updates(map[string]any{"created_at": createdAt, "updated_at": createdAt}).Error; err != nil {
		t.Fatalf("set chore timestamps: %v", err)
	}

	if completedDayOffset == 0 {
		return
	}

	completion := models.ChoreCompletion{ChoreID: chore.ID}
	if err := database.Create(&completion).Error; err != nil {
		t.Fatalf("seed chore completion: %v", err)
	}
	completedAt := day(completedDayOffset)
	if err := database.Model(&completion).Updates(map[string]any{"created_at": completedAt, "updated_at": completedAt}).Error; err != nil {
		t.Fatalf("set completion timestamps: %v", err)
	}
}

func fixedNow() time.Time {
	return day(0)
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(&models.Chore{}, &models.ChoreCompletion{})
}
