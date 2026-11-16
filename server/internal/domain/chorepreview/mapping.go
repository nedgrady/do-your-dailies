package chorepreview

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chorequeue"
)

func toAPIChorePreview(schedule []DaySchedule) contracts.ChorePreview {
	days := make([]contracts.ChorePreviewDay, 0, len(schedule))
	for _, day := range schedule {
		days = append(days, contracts.ChorePreviewDay{
			DayOffset: day.DayOffset,
			Chores:    toAPIChorePreviewChores(day.Chores),
		})
	}
	return contracts.ChorePreview{Days: days}
}

func toAPIChorePreviewChores(chores []chorequeue.ChoreInQueue) []contracts.ChorePreviewChore {
	result := make([]contracts.ChorePreviewChore, 0, len(chores))
	for _, chore := range chores {
		result = append(result, contracts.ChorePreviewChore{
			ChoreId:   uint64(chore.ChoreID),
			ChoreName: chore.ChoreName,
		})
	}
	return result
}
