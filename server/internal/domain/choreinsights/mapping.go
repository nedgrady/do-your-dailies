package choreinsights

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chores"
)

func toAPIChoreProjections(projections []ChoreProjection) []contracts.ChoreProjection {
	result := make([]contracts.ChoreProjection, 0, len(projections))
	for _, projection := range projections {
		result = append(result, toAPIChoreProjection(projection))
	}
	return result
}

func toAPIChoreProjection(projection ChoreProjection) contracts.ChoreProjection {
	return contracts.ChoreProjection{
		Chore:            chores.ToAPIChore(projection.Chore),
		ProjectedCadence: projection.ProjectedCadence,
	}
}
