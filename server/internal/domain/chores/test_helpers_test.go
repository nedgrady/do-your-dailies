package chores

import (
	"errors"
	"testing"

	"do-your-dailies/server/internal/domain/models"
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

func seedChore(t *testing.T, database *gorm.DB, name string, cadenceInDays int) models.Chore {
	t.Helper()

	chore := models.Chore{Name: name, CadenceInDays: cadenceInDays}
	if err := database.Create(&chore).Error; err != nil {
		t.Fatalf("seed models.Chore: %v", err)
	}

	return chore
}

type failingStore struct {
	listErr   error
	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func (store failingStore) List() ([]models.Chore, error) {
	if store.listErr != nil {
		return nil, store.listErr
	}
	return nil, errors.New("list not implemented")
}

func (store failingStore) Create(req CreateRequest) (models.Chore, error) {
	if store.createErr != nil {
		return models.Chore{}, store.createErr
	}
	return models.Chore{}, errors.New("create not implemented")
}

func (store failingStore) Get(id uint) (models.Chore, error) {
	if store.getErr != nil {
		return models.Chore{}, store.getErr
	}
	return models.Chore{}, errors.New("get not implemented")
}

func (store failingStore) Update(id uint, req UpdateRequest) (models.Chore, error) {
	if store.updateErr != nil {
		return models.Chore{}, store.updateErr
	}
	return models.Chore{}, errors.New("update not implemented")
}

func (store failingStore) Delete(id uint) error {
	if store.deleteErr != nil {
		return store.deleteErr
	}
	return errors.New("delete not implemented")
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
	return database.AutoMigrate(&models.Chore{})
}
