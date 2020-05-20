package services

import (
	"fitness/database"
	//"fitness/services/exercises"
	//"fitness/services/sessions"
	"fitness/services/users"
	//"fitness/services/workouts"
)

// Services defines the Services
type Services struct {
	Users     *users.Services
	Sessions  *sessions.Services
	Workouts  *workouts.Services
	Exercises *exercises.Services
}

// New returns a new Services.
func New(db *database.Database) *Services {
	return &Services{
		Users:     users.New(db),
		Sessions:  sessions.New(db),
		Workouts:  workouts.New(db),
		Exercises: exercises.New(db),
	}
}
