package models

import "gorm.io/gorm"

// ChoreCompletor is a named person on an account who can complete chores.
// It has no login of its own; UserID is the account it belongs to.
type ChoreCompletor struct {
	gorm.Model
	UserID uint   `json:"user_id" gorm:"not null;index"`
	Name   string `json:"name" gorm:"not null"`
}
