package workouts

import "errors"

var (
	// ErrLogNotFound is returned when a todo could not be found.
	ErrLogNotFound = errors.New("Workout log could not be found")
)
