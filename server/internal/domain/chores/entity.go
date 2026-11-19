package chores

import (
	"time"

	"do-your-dailies/server/internal/domain/models"
)

type CreateRequest struct {
	Name            string             `json:"name"`
	CadenceInDays   int                `json:"cadence_in_days"`
	DisplayUnit     models.DisplayUnit `json:"display_unit"`
	Emoji           models.Emoji       `json:"emoji"`
	LastCompletedAt *time.Time         `json:"last_completed_at"`
}

type UpdateRequest struct {
	Name          *string             `json:"name"`
	CadenceInDays *int                `json:"cadence_in_days"`
	DisplayUnit   *models.DisplayUnit `json:"display_unit"`
	Emoji         *models.Emoji       `json:"emoji"`
}
