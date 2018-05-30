package web

import (
	"fmt"
)

var (
	errBadRequest   = InnError{Code: 101, Message: "Bad Input Request"}
	errSlugNotMatch = InnError{Code: 102, Message: "Slug Not Match"}
	errSlugNotFound = InnError{Code: 103, Message: "Slug Not Found"}

	// Internal Server Error
	errInternalServer = InnError{Code: 500, Message: "Internal Server Error"}
)

// InnError is an error implementation that includes a time and message.
type InnError struct {
	Code    int
	Message string
}

func (e InnError) Error() string {
	return fmt.Sprintf("Error Code: %d, Error Message: %s", e.Code, e.Message)
}

// IsInnError is check error type
func IsInnError(err error) bool {
	_, ok := err.(InnError)
	return ok
}
