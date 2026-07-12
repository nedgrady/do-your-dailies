package models

import "gorm.io/gorm"

type ChoreInQueueCompletion struct {
	gorm.Model
	ChoreCompletionID uint `json:"chore_completion_id" gorm:"not null"`
}
