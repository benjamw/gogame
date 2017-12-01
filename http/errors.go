package http

import (
	"fmt"
	"net/http"
)

// MissingRequiredError gets thrown when a required form field is missing
type MissingRequiredError struct {
	FormElement string
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *MissingRequiredError) Error() string {
	return fmt.Sprintf("Missing required field: %s", e.FormElement)
}

// Code allows the struct to implement the game.Error interface
func (e *MissingRequiredError) Code() int {
	return http.StatusPreconditionFailed
}

// InvalidMethodError gets thrown when a request has an invalid method (GET vs POST, etc)
type InvalidMethodError struct {
	Method string
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *InvalidMethodError) Error() string {
	return fmt.Sprintf("Invalid Method: %s", e.Method)
}

// Code allows the struct to implement the game.Error interface
func (e *InvalidMethodError) Code() int {
	return http.StatusMethodNotAllowed
}

// BadLoginError gets thrown when the authentication for a user is invalid
type BadLoginError struct {
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *BadLoginError) Error() string {
	return "Invalid email and/or password"
}

// Code allows the struct to implement the game.Error interface
func (e *BadLoginError) Code() int {
	return http.StatusUnauthorized
}

// NotAuthorizedError gets thrown when an endpoint gets hit with invalid credentials
type NotAuthorizedError struct {
	Err error
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *NotAuthorizedError) Error() string {
	return "Not authorized to access this area"
}

// Code allows the struct to implement the game.Error interface
func (e *NotAuthorizedError) Code() int {
	return http.StatusUnauthorized
}

// RedirectError gets thrown when a URL or other path needs to redirect
// it's not necessarily an error
type RedirectError struct {
	URL    string
	Status int
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *RedirectError) Error() string {
	return e.URL
}

// Code allows the struct to implement the game.Error interface
func (e *RedirectError) Code() int {
	return e.Status
}
