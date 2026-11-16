package models

import "gorm.io/gorm"

type ChoreCompletion struct {
	gorm.Model
	ChoreID    uint              `json:"chore_id" gorm:"not null"`
	Chore      *Chore            `json:"chore" gorm:"foreignKey:ChoreID"`
	Completors []*ChoreCompletor `json:"completors" gorm:"many2many:chore_completion_completors"`
}
