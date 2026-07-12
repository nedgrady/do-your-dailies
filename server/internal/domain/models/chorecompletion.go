package models

import "gorm.io/gorm"

type ChoreCompletion struct {
	gorm.Model
	ChoreID uint   `json:"chore_id" gorm:"not null"`
	Chore   *Chore `json:"chore" gorm:"foreignKey:ChoreID"`
}
