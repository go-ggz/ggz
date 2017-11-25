package models

import (
	"fmt"
)

// ErrURLExist represents a "ErrURLExist" kind of error.
type ErrURLExist struct {
	Slug string
	URL  string
}

// IsErrURLExist checks if an error is a ErrURLExist.
func IsErrURLExist(err error) bool {
	_, ok := err.(ErrURLExist)
	return ok
}

func (err ErrURLExist) Error() string {
	return fmt.Sprintf("URL exist, slug: %s, url: %s", err.Slug, err.URL)
}

// ErrUserNotExist represents a "UserNotExist" kind of error.
type ErrUserNotExist struct {
	UID   int64
	Name  string
	KeyID int64
}

// IsErrUserNotExist checks if an error is a ErrUserNotExist.
func IsErrUserNotExist(err error) bool {
	_, ok := err.(ErrUserNotExist)
	return ok
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [uid: %d, name: %s, keyid: %d]", err.UID, err.Name, err.KeyID)
}
