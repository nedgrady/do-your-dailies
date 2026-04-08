package models

import "gorm.io/gorm"

type Chore struct {
	gorm.Model
	Name          string            `json:"name" gorm:"not null"`
	CadenceInDays int               `json:"cadence_in_days" gorm:"not null"`
	Completions   []ChoreCompletion `json:"completions"`
}

type ChoreCompletion struct {
	gorm.Model
	ChoreID uint  `json:"chore_id" gorm:"not null"`
	Chore   Chore `json:"chore"`
}

type CreateChoreRequest struct {
	Name          string `json:"name"`
	CadenceInDays int    `json:"cadence_in_days"`
}

type CreateChoreCompletionRequest struct {
	ChoreID uint `json:"chore_id"`
}

type UpdateChoreRequest struct {
	Name          *string `json:"name"`
	CadenceInDays *int    `json:"cadence_in_days"`
}
