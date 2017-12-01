package game

import (
	"fmt"
	"net/http"
)

// MultipleObjectError gets thrown when multiple entities are found when only one should exist
type MultipleObjectError struct {
	EntityType string // model.EntityType() response ("Company", "Asset", etc)
	Key        string // the entity key
	Value      string // the entity key value
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *MultipleObjectError) Error() string {
	return fmt.Sprintf("Multiple %ss found with %s of '%s'", e.EntityType, e.Key, e.Value)
}

// Code allows the struct to implement the game.Error interface
func (e *MultipleObjectError) Code() int {
	return http.StatusBadRequest
}

// DuplicateObjectError gets thrown when a duplicate entity already exists in the datastore
type DuplicateObjectError struct {
	EntityType string // model.EntityType() response ("Company", "Asset", etc)
	Key        string // the entity key
	Value      string // the entity key value
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *DuplicateObjectError) Error() string {
	return fmt.Sprintf("%s already exists with %s of '%s'", e.EntityType, e.Key, e.Value)
}

// Code allows the struct to implement the game.Error interface
func (e *DuplicateObjectError) Code() int {
	return http.StatusBadRequest
}

// AccountExistsError gets thrown when an account with the given email address already exists
type AccountExistsError struct {
	Email string // the email address used
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *AccountExistsError) Error() string {
	return fmt.Sprint("An account with the given credentials already exists on this system.")
}

// Code allows the struct to implement the game.Error interface
func (e *AccountExistsError) Code() int {
	return http.StatusBadRequest
}

// InvalidCredentialsError gets thrown when someone tries to log in with invalid credentials
// (missing account, wrong password, etc)
type InvalidCredentialsError struct {
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *InvalidCredentialsError) Error() string {
	return fmt.Sprint("The credentials supplied are invalid.")
}

// Code allows the struct to implement the game.Error interface
func (e *InvalidCredentialsError) Code() int {
	return http.StatusBadRequest
}

// UserError is an error type that is strictly used for error output to the end user.
type UserError struct {
	Status  int    // the http status code
	Message string // the error message to return
	Err     error  // the original error (if any)
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *UserError) Error() string {
	return e.Message
}

// Code allows the struct to implement the game.Error interface
func (e *UserError) Code() int {
	return e.Status
}

// NewUserError creates a UserError with the given format and parameters
func NewUserError(err error, code int, format string, params ...interface{}) error {
	return &UserError{
		Status:  code,
		Message: fmt.Sprintf(format, params...),
		Err:     err,
	}
}
