package chores

import (
	"errors"
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

type failingStore struct {
	listErr   error
	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func (store failingStore) List() ([]Chore, error) {
	if store.listErr != nil {
		return nil, store.listErr
	}
	return nil, errors.New("list not implemented")
}

func (store failingStore) Create(req CreateRequest) (Chore, error) {
	if store.createErr != nil {
		return Chore{}, store.createErr
	}
	return Chore{}, errors.New("create not implemented")
}

func (store failingStore) Get(id uint) (Chore, error) {
	if store.getErr != nil {
		return Chore{}, store.getErr
	}
	return Chore{}, errors.New("get not implemented")
}

func (store failingStore) Update(id uint, req UpdateRequest) (Chore, error) {
	if store.updateErr != nil {
		return Chore{}, store.updateErr
	}
	return Chore{}, errors.New("update not implemented")
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
	return database.AutoMigrate(&Chore{})
}
