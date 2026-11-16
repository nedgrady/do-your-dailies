package chorepreview

import (
	"context"

	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"

	"gorm.io/gorm"
)

// maxPreviewDays mirrors the OpenAPI spec's `maximum: 90` on the `days`
// query parameter - re-checked here defensively in case that validation
// is ever loosened without this handler being revisited.
const maxPreviewDays = 90

type ChorePreviewHandler struct {
	Db *gorm.DB
}

func NewHandler(db *gorm.DB) ChorePreviewHandler {
	return ChorePreviewHandler{Db: db}
}

func (handler ChorePreviewHandler) GetChorePreview(ctx context.Context, request contracts.GetChorePreviewRequestObject) (contracts.GetChorePreviewResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	days := request.Params.Days
	if days <= 0 || days > maxPreviewDays {
		return contracts.GetChorePreview400Response{}, nil
	}

	schedule, err := Simulate(handler.Db, ctx, userID, days)
	if err != nil {
		return nil, err
	}

	return contracts.GetChorePreview200JSONResponse(toAPIChorePreview(schedule)), nil
}
