package chorequeue

import (
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/chores"

	"github.com/stretchr/testify/assert"
)

func TestBuildQueueExcludesRecentlyCompletedLongCadenceChore(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "bathroom", 30, -100, intPtr(-5)),
		newCandidate(2, "dishes", 1, -100, intPtr(-1)),
		newCandidate(3, "vacuum", 7, -100, intPtr(-8)),
	}, day(0), 10)

	assert.Equal(t, []string{"vacuum", "dishes"}, queueNames(queue))
}

func TestBuildQueueOrdersByOverdueScore(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "bathroom", 30, -100, intPtr(-40)),
		newCandidate(2, "dishes", 1, -100, intPtr(-1)),
		newCandidate(3, "vacuum", 7, -100, intPtr(-8)),
	}, day(0), 10)

	assert.Equal(t, []string{"bathroom", "vacuum", "dishes"}, queueNames(queue))
}

func TestBuildQueueAppliesDailyCap(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "bathroom", 30, -100, intPtr(-40)),
		newCandidate(2, "dishes", 1, -100, intPtr(-1)),
		newCandidate(3, "vacuum", 7, -100, intPtr(-8)),
		newCandidate(4, "bins", 7, -100, intPtr(-20)),
	}, day(0), 2)

	assert.Equal(t, []string{"bins", "bathroom"}, queueNames(queue))
}

func TestBuildQueueUsesCreatedDayForNeverCompletedChores(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "wipe skirting", 14, -20, nil),
		newCandidate(2, "clean oven", 30, -20, nil),
		newCandidate(3, "dishes", 1, -100, intPtr(-1)),
	}, day(0), 10)

	assert.Equal(t, []string{"wipe skirting", "dishes"}, queueNames(queue))
}

func TestBuildQueueUsesDueDayTieBreakerWhenScoresMatch(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "desk tidy", 2, -20, intPtr(-3)),
		newCandidate(2, "laundry", 4, -20, intPtr(-6)),
	}, day(0), 10)

	assert.Equal(t, []string{"laundry", "desk tidy"}, queueNames(queue))
}

func TestBuildQueueReturnsEmptyWhenNothingIsDue(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "bathroom", 30, -100, intPtr(-1)),
		newCandidate(2, "windows", 60, -100, intPtr(-2)),
		newCandidate(3, "deep mop", 14, -5, nil),
	}, day(0), 10)

	assert.Empty(t, queue)
}

func newCandidate(id uint, name string, cadenceInDays int, createdDayOffset int, completedDayOffset *int) Candidate {
	chore := chores.Chore{
		Name:          name,
		CadenceInDays: cadenceInDays,
	}
	chore.ID = id
	chore.CreatedAt = day(createdDayOffset)

	var lastCompletedAt *time.Time
	if completedDayOffset != nil {
		completedAt := day(*completedDayOffset)
		lastCompletedAt = &completedAt
	}

	return Candidate{Chore: chore, LastCompletedAt: lastCompletedAt}
}

func queueNames(queue []chores.Chore) []string {
	names := make([]string, 0, len(queue))
	for _, chore := range queue {
		names = append(names, chore.Name)
	}

	return names
}

func day(offset int) time.Time {
	base := time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC)
	return base.AddDate(0, 0, offset)
}

func intPtr(value int) *int {
	return &value
}
