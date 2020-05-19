package sessions

import "errors"

var (
	// ErrSessionNotFound is returned when a Session could not be found.
	ErrSessionNotFound = errors.New("Session could not be found")
)
