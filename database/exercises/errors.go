package exercises

import "errors"

var (
	// ErrExerciseNotFound is returned when a User could not be found.
	ErrExerciseNotFound = errors.New("Exercise could not be found")
)
