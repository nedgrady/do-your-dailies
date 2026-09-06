package chorequeue

import (
	"context"
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func migrateForChoreQueueTests(database *gorm.DB) error {
	return database.AutoMigrate(&models.User{}, &models.Chore{}, &models.ChoreCompletion{})
}

func seedUser(t *testing.T, database *gorm.DB) uint {
	t.Helper()
	user := models.User{Name: "Test User"}
	require.NoError(t, database.Create(&user).Error)
	return user.ID
}

// A chore backdated well past its cadence must show up as overdue (priority > 1)
// immediately, not read as "just done now" via the no-completion-row fallback.
func TestListForCapacityFirstUserSurfacesBackdatedBacklogImmediately(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoreQueueTests)
	userID := seedUser(t, database)
	choreStore := chores.NewGormStore(database)
	queueStore := NewGormStore(database)

	lastCompletedAt := time.Now().AddDate(0, 0, -30)
	_, err := choreStore.Create(context.Background(), userID, chores.CreateRequest{
		Name:            "Neglected chore",
		CadenceInDays:   7,
		DisplayUnit:     models.DisplayUnitDay,
		LastCompletedAt: &lastCompletedAt,
	})
	require.NoError(t, err)

	items, err := queueStore.ListForCapacityFirstUser(context.Background(), userID, 5)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Greater(t, items[0].Priority, 1.0)
}

// A chore created with no last-done info still defaults to "just done now" -
// today's existing fallback behavior, preserved for chores with no known history.
func TestListForCapacityFirstUserTreatsFreshChoreAsJustDone(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoreQueueTests)
	userID := seedUser(t, database)
	choreStore := chores.NewGormStore(database)
	queueStore := NewGormStore(database)

	_, err := choreStore.Create(context.Background(), userID, chores.CreateRequest{
		Name:          "Fresh chore",
		CadenceInDays: 7,
		DisplayUnit:   models.DisplayUnitDay,
	})
	require.NoError(t, err)

	items, err := queueStore.ListForCapacityFirstUser(context.Background(), userID, 5)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Less(t, items[0].Priority, 0.01)
}
