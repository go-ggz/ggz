package models

import (
	"time"
)

// User represents the object of individual and member of organization.
type User struct {
	ID       int64 `xorm:"pk autoincr"`
	FullName string
	// Email is the primary email address (to be used for communication)
	Email     string `xorm:"NOT NULL"`
	Passwd    string `xorm:"NOT NULL"`
	Location  string
	Website   string
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin time.Time
}
