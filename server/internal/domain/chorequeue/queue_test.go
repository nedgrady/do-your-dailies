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

func TestStartOfDayUTCReturnsMidnight(t *testing.T) {
	t.Parallel()

	dt := time.Date(2026, time.April, 15, 13, 14, 15, 123456789, time.FixedZone("X", 2*3600))
	got := startOfDayUTC(dt)
	want := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, want, got)
}

func TestBuildQueueIncludesOneDayOverdue(t *testing.T) {
	t.Parallel()

	anchor := day(0)
	chore := chores.Chore{Name: "oneDay", CadenceInDays: 7}
	chore.CreatedAt = anchor
	candidate := Candidate{Chore: chore, LastCompletedAt: nil}

	queue := BuildQueue([]Candidate{candidate}, day(8), 10)

	assert.Equal(t, []string{"oneDay"}, queueNames(queue))
}

func TestBuildQueueExcludesZeroCadence(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "zero", 0, -1, nil),
		newCandidate(2, "valid", 1, -1, intPtr(-2)),
	}, day(0), 10)

	assert.Equal(t, []string{"valid"}, queueNames(queue))
}

func TestBuildQueueOverdueScorePreferredToDueDay(t *testing.T) {
	t.Parallel()

	// Create two chores that share the same dueDay but have different cadences
	// so their overdueScore differs; the higher overdueScore should come first.
	// Candidate A: ID 1, created day 0, cadence 2 -> dueDay = day(2)
	// Candidate B: ID 2, created day 1, cadence 1 -> dueDay = day(2)
	a := chores.Chore{Name: "A", CadenceInDays: 2}
	a.ID = 1
	a.CreatedAt = day(0)
	b := chores.Chore{Name: "B", CadenceInDays: 1}
	b.ID = 2
	b.CreatedAt = day(1)

	queue := BuildQueue([]Candidate{
		{Chore: a, LastCompletedAt: nil},
		{Chore: b, LastCompletedAt: nil},
	}, day(3), 10)

	assert.Equal(t, []string{"B", "A"}, queueNames(queue))
}

func TestBuildQueueNoCapWhenMaxChoresNonPositive(t *testing.T) {
	t.Parallel()

	queue := BuildQueue([]Candidate{
		newCandidate(1, "one", 1, -10, intPtr(-2)),
		newCandidate(2, "two", 1, -10, intPtr(-2)),
	}, day(0), 0)

	assert.Equal(t, 2, len(queue))
}

func TestBuildQueueAppliesStartOfDayToTarget(t *testing.T) {
	t.Parallel()

	chore := chores.Chore{Name: "X", CadenceInDays: 1}
	chore.CreatedAt = day(0)
	candidate := Candidate{Chore: chore}

	q1 := BuildQueue([]Candidate{candidate}, time.Date(2026, time.April, 16, 12, 0, 0, 0, time.UTC), 10)
	q2 := BuildQueue([]Candidate{candidate}, time.Date(2026, time.April, 16, 0, 0, 0, 0, time.UTC), 10)

	assert.Equal(t, q1, q2)
}

func TestBuildQueueDueDayTieBreakerUsedWhenScoresEqual(t *testing.T) {
	t.Parallel()

	// Construct two chores with equal overdueScore but different dueDays.
	// Expect the earlier dueDay to come first.
	a := chores.Chore{Name: "A", CadenceInDays: 4}
	a.ID = 1
	a.CreatedAt = day(0) // dueDay = day(4)

	b := chores.Chore{Name: "B", CadenceInDays: 2}
	b.ID = 2
	b.CreatedAt = day(3) // dueDay = day(5)

	queue := BuildQueue([]Candidate{{Chore: b}, {Chore: a}}, day(6), 10)

	assert.Equal(t, []string{"A", "B"}, queueNames(queue))
}
