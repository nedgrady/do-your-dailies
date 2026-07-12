package chorequeue

import (
	"testing"

	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestListForCapacityFirstUserReturnsChoresOrderedByCompletionRatio(t *testing.T) {
	t.Parallel()

	database := testhelpers.NewTransactionDB(t, migrate)
	store := NewGormStore(database)

	seedChoreWithCompletion(t, database, "bathroom", 30, -100, -40)
	seedChoreWithCompletion(t, database, "dishes", 1, -100, -1)

	chores, err := store.ListForCapacityFirstUser(2)

	assert.NoError(t, err)
	assert.Len(t, chores, 2)
	assert.Equal(t, []string{"bathroom", "dishes"}, choreNames(chores))
}

func choreNames(chores []ChoreInQueue) []string {
	names := make([]string, 0, len(chores))
	for _, chore := range chores {
		names = append(names, chore.ChoreName)
	}

	return names
}
