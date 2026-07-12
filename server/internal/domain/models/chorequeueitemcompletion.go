package models

import "gorm.io/gorm"

type ChoreQueueItemCompletion struct {
	gorm.Model
	ChoreCompletionId uint `json:"chore_id" gorm:"not null"`
}
