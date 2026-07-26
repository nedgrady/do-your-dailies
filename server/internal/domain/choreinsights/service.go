package choreinsights

import (
	"do-your-dailies/server/internal/domain/models"
	"fmt"
)

const defaultUserDesiredCapacity = 5

type Insights struct {
	UserDesiredCapacity                             int
	UtilizationRatio                                float64
	MinCapacityToKeepUtilizationRatioGreaterThanOne float64
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Calculate(chores []models.Chore) (Insights, error) {
	var totalFrequency float64
	for _, chore := range chores {
		if chore.CadenceInDays <= 0 {
			return Insights{}, fmt.Errorf("chore %d has invalid cadenceInDays: %d", chore.ID, chore.CadenceInDays)
		}
		totalFrequency += 1.0 / float64(chore.CadenceInDays)
	}

	// TODO: replace with a per-user setting once one exists; hardcoded for now.
	userDesiredCapacity := defaultUserDesiredCapacity
	if userDesiredCapacity <= 0 {
		return Insights{}, fmt.Errorf("userDesiredCapacity must be greater than zero, got: %d", userDesiredCapacity)
	}

	utilizationRatio := totalFrequency / float64(userDesiredCapacity)

	return Insights{
		UserDesiredCapacity: userDesiredCapacity,
		UtilizationRatio:    utilizationRatio,
		MinCapacityToKeepUtilizationRatioGreaterThanOne: totalFrequency,
	}, nil
}
