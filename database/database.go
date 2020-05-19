package database

import (
	"database/sql"
	"fitness/database/exercises"
	"fitness/database/sessions"
	"fitness/database/users"
	"fitness/database/workouts"
)

// Database defines the database
type Database struct {
	Users     *users.Database
	Sessions  *sessions.Database
	Workouts  *workouts.Database
	Exercises *exercises.Database
}

// New returns a new database.
func New(db *sql.DB) *Database {
	return &Database{
		Users:     users.New(db),
		Sessions:  sessions.New(db),
		Workouts:  workouts.New(db),
		Exercises: exercises.New(db),
	}
}
