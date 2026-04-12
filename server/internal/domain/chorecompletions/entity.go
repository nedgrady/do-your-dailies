package chorecompletions

import "gorm.io/gorm"

type ChoreCompletion struct {
	gorm.Model
	ChoreID uint `json:"chore_id" gorm:"not null"`
}

type CreateRequest struct {
	ChoreID uint `json:"chore_id"`
}
