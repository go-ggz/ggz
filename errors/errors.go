package errors

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Error applicational
type Error struct {
	Type    Type
	Message string
	cause   error
}

// Error message
func (e *Error) Error() string {
	return e.Message
}

// Cause of the original error
func (e *Error) Cause() string {
	if e.cause != nil {
		return e.cause.Error()
	}

	return ""
}

// Extensions for graphQL extension
func (e *Error) Extensions() map[string]interface{} {
	if e.cause != nil {
		log.Error().Err(e.cause).Msg("graphql error report")
	}

	return map[string]interface{}{
		"code": e.Type.Code(),
		"type": e.Type,
	}
}

// Type defines the type of an error
type Type string

const (
	// Internal error
	Internal Type = "internal"
	// NotFound error means that a specific item does not exist
	NotFound Type = "not_found"
	// BadRequest error
	BadRequest Type = "bad_request"
	// Validation error
	Validation Type = "validation"
	// AlreadyExists error
	AlreadyExists Type = "already_exists"
	// Unauthorized error
	Unauthorized Type = "unauthorized"
)

func (t Type) String() string {
	switch t {
	case Internal:
		return "Internal Error"
	case NotFound:
		return "Item not found"
	case BadRequest:
		return "BadRequest error"
	case Validation:
		return "Validation error"
	case AlreadyExists:
		return "Item already exists"
	case Unauthorized:
		return "Unauthorized error"
	}

	return "Unknown error"
}

// Code http error code
func (t Type) Code() int {
	switch t {
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case BadRequest:
		return http.StatusBadRequest
	case Validation:
		return http.StatusForbidden
	case AlreadyExists:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

// New creates a new error
func New(t Type, msg string, err error) error {
	return &Error{
		Type:    t,
		Message: msg,
		cause:   err,
	}
}

// EValidation creates an error of type Validationn
func EValidation(msg string, err error, arg ...interface{}) error {
	return New(Validation, fmt.Sprintf(msg, arg...), err)
}

// ENotExists creates an error of type NotExist
func ENotExists(msg string, err error, arg ...interface{}) error {
	return New(NotFound, fmt.Sprintf(msg, arg...), err)
}

// EBadRequest creates an error of type BadRequest
func EBadRequest(msg string, err error, arg ...interface{}) error {
	return New(BadRequest, fmt.Sprintf(msg, arg...), err)
}

// EAlreadyExists creates an error of type EAlreadyExistsl
func EAlreadyExists(msg string, err error, arg ...interface{}) error {
	return New(AlreadyExists, fmt.Sprintf(msg, arg...), err)
}

// EInternal creates an error of type Internal
func EInternal(msg string, err error, arg ...interface{}) error {
	return New(Internal, fmt.Sprintf(msg, arg...), err)
}

// ENotFound creates an error of type NotFound
func ENotFound(msg string, err error, arg ...interface{}) error {
	return New(NotFound, fmt.Sprintf(msg, arg...), err)
}

// EUnauthorized creates an error of type Unauthorized
func EUnauthorized(msg string, err error, arg ...interface{}) error {
	return New(Unauthorized, fmt.Sprintf(msg, arg...), err)
}

// Is method checks if an error is of a specific type
func Is(t Type, err error) bool {
	e, ok := err.(*Error)

	return ok && e.Type == t
}
