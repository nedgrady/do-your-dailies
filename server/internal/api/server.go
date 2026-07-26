package api

import (
	"context"
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/choreinsights"
	"do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
)

type Server struct {
	chores.ChoreHandler
	chorecompletions.ChoreCompletionHandler
	chorequeue.ChoreQueueHandler
	choreinsights.ChoreInsightsHandler
}

func (Server) HealthCheck(ctx context.Context, request contracts.HealthCheckRequestObject) (contracts.HealthCheckResponseObject, error) {
	return contracts.HealthCheck200TextResponse("OK"), nil
}
