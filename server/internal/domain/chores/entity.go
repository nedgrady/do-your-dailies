package chores

import "gorm.io/gorm"

type Chore struct {
	gorm.Model
	Name          string `json:"name" gorm:"not null"`
	CadenceInDays int    `json:"cadence_in_days" gorm:"not null"`
}

type CreateRequest struct {
	Name          string `json:"name"`
	CadenceInDays int    `json:"cadence_in_days"`
}

type UpdateRequest struct {
	Name          *string `json:"name"`
	CadenceInDays *int    `json:"cadence_in_days"`
}
