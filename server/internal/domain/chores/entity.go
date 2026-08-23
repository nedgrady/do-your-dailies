package chores

import "do-your-dailies/server/internal/domain/models"

type CreateRequest struct {
	Name          string             `json:"name"`
	CadenceInDays int                `json:"cadence_in_days"`
	DisplayUnit   models.DisplayUnit `json:"display_unit"`
}

type UpdateRequest struct {
	Name          *string             `json:"name"`
	CadenceInDays *int                `json:"cadence_in_days"`
	DisplayUnit   *models.DisplayUnit `json:"display_unit"`
}
