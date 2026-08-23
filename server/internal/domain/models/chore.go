package models

import "gorm.io/gorm"

type Chore struct {
	gorm.Model
	Name             string             `json:"name" gorm:"not null"`
	CadenceInDays    int                `json:"cadence_in_days" gorm:"not null"`
	DisplayUnit      DisplayUnit        `json:"display_unit" gorm:"not null;default:'DAY'"`
	ChoreCompletions []*ChoreCompletion `json:"chore_completions" gorm:"foreignKey:ChoreID"`
	UserID           uint               `json:"user_id" gorm:"not null"`
}
