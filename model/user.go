package model

import (
	"strings"
	"time"

	"github.com/go-ggz/ggz/module/base"
)

// User represents the object of individual and member of organization.
type User struct {
	ID       int64  `xorm:"pk autoincr" json:"id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	// Email is the primary email address (to be used for communication)
	Email       string `xorm:"UNIQUE NOT NULL" json:"email,omitempty"`
	Location    string
	Website     string
	IsActive    bool      `xorm:"INDEX"` // Activate primary email
	Avatar      string    `xorm:"VARCHAR(2048) NOT NULL" json:"avatar,omitempty"`
	AvatarEmail string    `xorm:"NOT NULL" json:"avatar_email,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	LastLogin   time.Time `json:"lastlogin,omitempty"`
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

// GetUserByEmail returns the user object by given e-mail if exists.
func GetUserByEmail(email string) (*User, error) {
	if len(email) == 0 {
		return nil, ErrUserNotExist{0, email, 0}
	}

	email = strings.ToLower(email)
	// First try to find the user by primary email
	user := &User{Email: email}
	has, err := x.Get(user)
	if err != nil {
		return nil, err
	}
	if has {
		return user, nil
	}

	return nil, ErrUserNotExist{0, email, 0}
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

	u.AvatarEmail = u.Email
	u.Avatar = base.HashEmail(u.AvatarEmail)

	if _, err = sess.Insert(u); err != nil {
		return err
	}

	return sess.Commit()
}
