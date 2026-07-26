package choreinsights

import (
	"context"

	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chores"
)

type ChoreInsightsHandler struct {
	ChoreStore chores.Store
	Service    Service
}

func NewHandler(choreStore chores.Store) ChoreInsightsHandler {
	return ChoreInsightsHandler{ChoreStore: choreStore, Service: NewService()}
}

func (handler ChoreInsightsHandler) GetChoreInsights(ctx context.Context, request contracts.GetChoreInsightsRequestObject) (contracts.GetChoreInsightsResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userChores, err := handler.ChoreStore.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	insights, err := handler.Service.Calculate(userChores)
	if err != nil {
		return nil, err
	}

	return contracts.GetChoreInsights200JSONResponse{
		UserDesiredCapacity: insights.UserDesiredCapacity,
		UtilizationRatio:    insights.UtilizationRatio,
		MinCapacityToKeepUtilizationRatioGreaterThanOne: insights.MinCapacityToKeepUtilizationRatioGreaterThanOne,
	}, nil
}
