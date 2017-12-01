package session

import (
	"net/http"
)

// InvalidSessionError is thrown when a session is invalid for whatever reason
type InvalidSessionError struct {
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *InvalidSessionError) Error() string {
	return "Invalid Session: Please log in again"
}

// Code allows the struct to implement the game.Error interface
func (e *InvalidSessionError) Code() int {
	return http.StatusUnauthorized
}

// ExpiredSessionError is thrown when a session has expired
type ExpiredSessionError struct {
}

// Error allows the struct to implement the error interface as well as the game.Error interface
func (e *ExpiredSessionError) Error() string {
	return "Session Expired: Please log in again"
}

// Code allows the struct to implement the game.Error interface
func (e *ExpiredSessionError) Code() int {
	return http.StatusUnauthorized
}
