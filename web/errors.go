package web

import (
	"fmt"
)

const (
	errInput = 100
	errDB    = 200
)

var (
	errBadRequest   = InnError{Code: 101, Message: "Bad Input Request"}
	errSlugNotMatch = InnError{Code: 102, Message: "Slug Not Match"}
	errSlugNotFound = InnError{Code: 103, Message: "Slug Not Found"}
	errURLExist     = InnError{Code: 104, Message: "URL Exist"}

	// Internal Server Error
	errInternalServer = InnError{Code: 500, Message: "Internal Server Error"}
	errDBQuery        = InnError{Code: 501, Message: "Internal Database Error"}
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
