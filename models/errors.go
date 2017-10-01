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
