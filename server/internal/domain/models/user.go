package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name   string  `json:"name" gorm:"not null"`
	Chores []Chore `json:"chores" gorm:"foreignKey:UserID"`
}
