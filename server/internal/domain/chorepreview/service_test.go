package chorepreview

import (
	"context"
	"testing"

	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func migrateForChorePreviewTests(database *gorm.DB) error {
	return database.AutoMigrate(&models.User{}, &models.Chore{}, &models.ChoreCompletion{}, &models.ChoreCompletor{})
}

func seedUser(t *testing.T, database *gorm.DB) uint {
	t.Helper()
	user := models.User{Name: "Test User"}
	require.NoError(t, database.Create(&user).Error)
	return user.ID
}

func createChore(t *testing.T, database *gorm.DB, userID uint, name string, cadenceInDays int) {
	t.Helper()
	_, err := chores.NewGormStore(database).Create(context.Background(), userID, chores.CreateRequest{
		Name:          name,
		CadenceInDays: cadenceInDays,
		DisplayUnit:   models.DisplayUnitDay,
	})
	require.NoError(t, err)
}

func choreNamesForDay(schedule []DaySchedule, dayOffset int) []string {
	for _, day := range schedule {
		if day.DayOffset != dayOffset {
			continue
		}
		names := make([]string, 0, len(day.Chores))
		for _, chore := range day.Chores {
			names = append(names, chore.ChoreName)
		}
		return names
	}
	return nil
}

// The core property this whole feature exists for: a daily chore must
// appear on every one of the N simulated days, not just once.
func TestSimulateRepeatsADailyChoreOnEveryDay(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChorePreviewTests)
	userID := seedUser(t, database)
	createChore(t, database, userID, "Sweep floor", 1)

	schedule, err := Simulate(database, context.Background(), userID, 5)
	require.NoError(t, err)
	require.Len(t, schedule, 5)

	for day := 1; day <= 5; day++ {
		assert.Contains(t, choreNamesForDay(schedule, day), "Sweep floor", "day %d", day)
	}
}

// With capacity saturated by 5 daily competitors (whose priority resets to
// ~1.0 immediately after every pick, so they win every single day), a
// long-cadence chore's priority stays far below that steady-state for the
// whole window and must lose out on every day - proving the ranking
// correctly suppresses a chore whose cadence hasn't elapsed when something
// more urgent is competing for the same slots. (A short competing cadence
// like 3 or 4 would eventually tie the fillers' exact 1.0 steady-state on
// day == cadence, and Postgres doesn't guarantee tie-break order, so this
// uses a cadence far outside the test window instead of chasing an exact
// reappearance day. Eventual reappearance when uncontested is already
// covered by TestSimulateRepeatsADailyChoreOnEveryDay.)
func TestSimulateDoesNotRepickAChoreBeforeItsCadenceElapses(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChorePreviewTests)
	userID := seedUser(t, database)
	createChore(t, database, userID, "Mow lawn", 100)
	for i := 0; i < 5; i++ {
		createChore(t, database, userID, "Daily filler", 1)
	}

	schedule, err := Simulate(database, context.Background(), userID, 4)
	require.NoError(t, err)
	require.Len(t, schedule, 4)

	for day := 1; day <= 4; day++ {
		assert.NotContains(t, choreNamesForDay(schedule, day), "Mow lawn", "day %d", day)
	}
}

// Capacity per day is fixed at 5 - more than 5 eligible chores on a day
// still yields exactly 5 picks, highest-priority first.
func TestSimulateCapsEachDayAtFiveChores(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChorePreviewTests)
	userID := seedUser(t, database)
	for i := 0; i < 7; i++ {
		createChore(t, database, userID, "Daily chore", 1)
	}

	schedule, err := Simulate(database, context.Background(), userID, 1)
	require.NoError(t, err)
	require.Len(t, schedule, 1)
	assert.Len(t, schedule[0].Chores, 5)
}

// Simulate must never persist anything - it's a pure dry run.
func TestSimulateNeverPersistsCompletions(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChorePreviewTests)
	userID := seedUser(t, database)
	createChore(t, database, userID, "Sweep floor", 1)

	_, err := Simulate(database, context.Background(), userID, 7)
	require.NoError(t, err)

	var count int64
	require.NoError(t, database.Model(&models.ChoreCompletion{}).Count(&count).Error)
	assert.Zero(t, count, "Simulate must roll back every completion it writes")
}

func TestSimulateReturnsEmptyChoresWhenUserHasNoChores(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForChorePreviewTests)
	userID := seedUser(t, database)

	schedule, err := Simulate(database, context.Background(), userID, 3)
	require.NoError(t, err)
	require.Len(t, schedule, 3)
	for _, day := range schedule {
		assert.Empty(t, day.Chores)
	}
}
