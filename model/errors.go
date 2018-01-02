package model

import "net/http"

// EditingExistingTokenError gets thrown when an attempt to edit an existing token is made
type EditingExistingTokenError struct {
}

func (e *EditingExistingTokenError) Error() string {
	return "Cannot edit an existing token"
}

func (e *EditingExistingTokenError) Code() int {
	return http.StatusUnauthorized
}
