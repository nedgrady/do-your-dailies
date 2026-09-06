package chores

import (
	"context"
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func migrateForChoresTests(database *gorm.DB) error {
	return database.AutoMigrate(&models.User{}, &models.Chore{}, &models.ChoreCompletion{})
}

func seedUser(t *testing.T, database *gorm.DB) uint {
	t.Helper()
	user := models.User{Name: "Test User"}
	require.NoError(t, database.Create(&user).Error)
	return user.ID
}

func TestCreateWithoutLastCompletedAtLeavesNoCompletionHistory(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoresTests)
	userID := seedUser(t, database)
	store := NewGormStore(database)

	chore, err := store.Create(context.Background(), userID, CreateRequest{
		Name:          "Water plants",
		CadenceInDays: 2,
		DisplayUnit:   models.DisplayUnitDay,
	})
	require.NoError(t, err)

	var count int64
	require.NoError(t, database.Model(&models.ChoreCompletion{}).Where("chore_id = ?", chore.ID).Count(&count).Error)
	assert.Zero(t, count)
}

func TestCreateWithLastCompletedAtBackdatesACompletion(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoresTests)
	userID := seedUser(t, database)
	store := NewGormStore(database)

	lastCompletedAt := time.Now().AddDate(0, 0, -30).Truncate(time.Second)
	chore, err := store.Create(context.Background(), userID, CreateRequest{
		Name:            "Mow lawn",
		CadenceInDays:   14,
		DisplayUnit:     models.DisplayUnitDay,
		LastCompletedAt: &lastCompletedAt,
	})
	require.NoError(t, err)

	var completions []models.ChoreCompletion
	require.NoError(t, database.Where("chore_id = ?", chore.ID).Find(&completions).Error)
	require.Len(t, completions, 1)
	assert.WithinDuration(t, lastCompletedAt, completions[0].CreatedAt, time.Second)
}

func TestCreateProvisionsAUserRowForANotYetSeenSyntheticUserID(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoresTests)
	store := NewGormStore(database)
	const syntheticUserID = 999

	chore, err := store.Create(context.Background(), syntheticUserID, CreateRequest{
		Name:          "Water plants",
		CadenceInDays: 2,
		DisplayUnit:   models.DisplayUnitDay,
	})
	require.NoError(t, err)
	assert.Equal(t, uint(syntheticUserID), chore.UserID)

	var user models.User
	require.NoError(t, database.First(&user, syntheticUserID).Error)
}

func TestCreateDoesNotDuplicateOrModifyAnAlreadyExistingUser(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoresTests)
	userID := seedUser(t, database)
	store := NewGormStore(database)

	_, err := store.Create(context.Background(), userID, CreateRequest{
		Name:          "Water plants",
		CadenceInDays: 2,
		DisplayUnit:   models.DisplayUnitDay,
	})
	require.NoError(t, err)

	var count int64
	require.NoError(t, database.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error)
	assert.Equal(t, int64(1), count)

	var user models.User
	require.NoError(t, database.First(&user, userID).Error)
	assert.Equal(t, "Test User", user.Name)
}

func TestCreateManyBackdatesCompletionsOnlyForRowsThatProvideThem(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChoresTests)
	userID := seedUser(t, database)
	store := NewGormStore(database)

	lastCompletedAt := time.Now().AddDate(0, 0, -10).Truncate(time.Second)
	created, err := store.CreateMany(context.Background(), userID, []CreateRequest{
		{Name: "Backdated", CadenceInDays: 7, DisplayUnit: models.DisplayUnitDay, LastCompletedAt: &lastCompletedAt},
		{Name: "Fresh", CadenceInDays: 7, DisplayUnit: models.DisplayUnitDay},
	})
	require.NoError(t, err)
	require.Len(t, created, 2)

	var backdatedCount int64
	require.NoError(t, database.Model(&models.ChoreCompletion{}).Where("chore_id = ?", created[0].ID).Count(&backdatedCount).Error)
	assert.Equal(t, int64(1), backdatedCount)

	var freshCount int64
	require.NoError(t, database.Model(&models.ChoreCompletion{}).Where("chore_id = ?", created[1].ID).Count(&freshCount).Error)
	assert.Zero(t, freshCount)
}
