package choreinqueuecompletion

import (
	"testing"

	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreatePersistsChoreInQueueCompletion(t *testing.T) {
	t.Parallel()

	database := testhelpers.NewTransactionDB(t, migrate)
	store := NewGormStore(database)

	created, err := store.Create(CreateRequest{ChoreCompletionID: 42})

	assert.NoError(t, err)
	assert.Equal(t, uint(42), created.ChoreCompletionID)
	assert.NotZero(t, created.ID)

	var persisted models.ChoreInQueueCompletion
	err = database.First(&persisted, created.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, uint(42), persisted.ChoreCompletionID)
}

func migrate(database *gorm.DB) error {
	return database.AutoMigrate(&models.ChoreInQueueCompletion{})
}
