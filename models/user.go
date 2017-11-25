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

func getUserByID(e Engine, id int64) (*User, error) {
	u := new(User)
	has, err := e.ID(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist{id, "", 0}
	}
	return u, nil
}

// GetUserByID returns the user object by given ID if exists.
func GetUserByID(id int64) (*User, error) {
	return getUserByID(x, id)
}
