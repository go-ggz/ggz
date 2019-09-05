package model

import (
	"errors"
	"fmt"
)

// ErrURLExist represents a "ErrURLExist" kind of error.
type ErrURLExist struct {
	Slug string
	URL  string
}

// IsErrURLExist checks if an error is a ErrURLExist.
func IsErrURLExist(err error) bool {
	return errors.As(err, &ErrURLExist{})
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
	return errors.As(err, &ErrUserNotExist{})
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [uid: %d, name: %s, keyid: %d]", err.UID, err.Name, err.KeyID)
}

// ErrShortenNotExist represents a "ShortenNotExist" kind of error.
type ErrShortenNotExist struct {
	Slug string
}

// IsErrShortenNotExist checks if an error is a ErrUserNotExist.
func IsErrShortenNotExist(err error) bool {
	return errors.As(err, &ErrShortenNotExist{})
}

func (err ErrShortenNotExist) Error() string {
	return fmt.Sprintf("shorten slug does not exist [slug: %s]", err.Slug)
}

// ErrUserAlreadyExist represents a "user already exists" error.
type ErrUserAlreadyExist struct {
	Name string
}

// IsErrUserAlreadyExist checks if an error is a ErrUserAlreadyExists.
func IsErrUserAlreadyExist(err error) bool {
	return errors.As(err, &ErrUserAlreadyExist{})
}

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("user already exists [name: %s]", err.Name)
}

// ErrEmailAlreadyUsed represents a "EmailAlreadyUsed" kind of error.
type ErrEmailAlreadyUsed struct {
	Email string
}

// IsErrEmailAlreadyUsed checks if an error is a ErrEmailAlreadyUsed.
func IsErrEmailAlreadyUsed(err error) bool {
	return errors.As(err, &ErrEmailAlreadyUsed{})
}

func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("e-mail has been used [email: %s]", err.Email)
}

// ErrAccessTokenNotExist represents a "AccessTokenNotExist" kind of error.
type ErrAccessTokenNotExist struct {
	SHA string
}

// IsErrAccessTokenNotExist checks if an error is a ErrAccessTokenNotExist.
func IsErrAccessTokenNotExist(err error) bool {
	return errors.As(err, &ErrAccessTokenNotExist{})
}

func (err ErrAccessTokenNotExist) Error() string {
	return fmt.Sprintf("access token does not exist [sha: %s]", err.SHA)
}

// ErrAccessTokenEmpty represents a "AccessTokenEmpty" kind of error.
type ErrAccessTokenEmpty struct {
}

// IsErrAccessTokenEmpty checks if an error is a ErrAccessTokenEmpty.
func IsErrAccessTokenEmpty(err error) bool {
	return errors.As(err, &ErrAccessTokenEmpty{})
}

func (err ErrAccessTokenEmpty) Error() string {
	return fmt.Sprintf("access token is empty")
}
