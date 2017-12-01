package models

import (
	"strings"
	"time"

	"github.com/go-ggz/ggz/modules/base"
)

// User represents the object of individual and member of organization.
type User struct {
	ID       int64 `xorm:"pk autoincr"`
	FullName string
	// Email is the primary email address (to be used for communication)
	Email       string `xorm:"UNIQUE NOT NULL" json:"email,omitempty"`
	UserName    string `xorm:"UNIQUE NULL" json:"username,omitempty"`
	Passwd      string `xorm:"NOT NULL"`
	Location    string
	Website     string
	IsActive    bool   `xorm:"INDEX"` // Activate primary email
	Avatar      string `xorm:"VARCHAR(2048) NOT NULL" json:"avatar,omitempty"`
	AvatarEmail string `xorm:"NOT NULL"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLogin   time.Time
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

func isUserExist(e Engine, uid int64, email string) (bool, error) {
	if len(email) == 0 {
		return false, nil
	}
	return e.
		Where("id!=?", uid).
		Get(&User{Email: strings.ToLower(email)})
}

// IsUserExist checks if given user email exist,
// the user email should be noncased unique.
// If uid is presented, then check will rule out that one,
// it is used when update a user email in settings page.
func IsUserExist(uid int64, email string) (bool, error) {
	return isUserExist(x, uid, email)
}

// CreateUser creates record of a new user.
func CreateUser(u *User) (err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	isExist, err := sess.
		Where("email=?", u.Email).
		Get(new(User))
	if err != nil {
		return err
	} else if isExist {
		return ErrEmailAlreadyUsed{u.Email}
	}

	if u.UserName != "" {
		isExist, err := isUserExist(sess, 0, u.UserName)
		if err != nil {
			return err
		} else if isExist {
			return ErrUserAlreadyExist{u.UserName}
		}
	}

	u.AvatarEmail = u.Email
	u.Avatar = base.HashEmail(u.AvatarEmail)

	if _, err = sess.Insert(u); err != nil {
		return err
	}

	return sess.Commit()
}
