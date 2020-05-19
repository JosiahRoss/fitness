package exercises

import (
	"database/sql"
	"fmt"
)

// Database defines the exercises database.
type Database struct {
	db *sql.DB
}

// New creates a new exercises database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Exercise defines a exercise
type Exercise struct {
	ID           int    `json:"id"`
	ExerciseName string `json:"exercise_name"`
	MuscleGroup  string `json:"muscle_group"`
	Description  string `json:"description"`
}

// Exercises defines a set of exercises.
type Exercises struct {
	Exercises []*Exercise `json:"exercises"`
	Total     int         `json:"total"`
}

// A of constants defines the SQL statments that will
// used for the exercises database.
const (

	// stmtSelectAll defines the SQl statement to
	// select all workouts in table
	stmtSelectAll = `
SELECT ALL
FROM exercises

`

	// stmtSelectCount defines the SQL statement to select
	// the total number of exercises, accouring to the filters.
	stmtSelectCount = `
SELECT COUNT(*)
FROM exercises
`

	// stmtSelectByID defines the SQL statement to select a
	// exercise by id.
	stmtSelectByID = `
SELECT id, exercise_name, muscle_group, description
FROM exercises
WEHRE id=?

`

	// stmtSelectByName defines the SQL statement to select a
	// exercise by name.
	stmtSelectByName = `
SELECT id, exercise_name, muscle_group, description
FROM exercises
WHERE name=?
`
	// stmtSelectByMuscleGroup defines the SQL statement to
	// select a group of exercises by muscle_group
	stmtSelectByMuscleGroup = `
SELECT 	id, exercise_name, muscle_group, description
FROM exercises
WHERE muscle_group=?
`
)

// GetAll returns all exercises
func (db *Database) GetAll() (*Exercises, error) {
	// Build the full query.
	query := fmt.Sprintf(stmtSelectAll)

	// Create a new Exercises.
	exercises := &Exercises{
		Exercises: []*Exercise{},
	}

	// Execute the querey.
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// loop through the exercise row.
	for rows.Next() {
		// Create a new Exercise
		exercise := &Exercise{}

		// Scan row values into Exercise struct.
		if err := rows.Scan(&exercise.ID, &exercise.ExerciseName, &exercise.MuscleGroup, &exercise.Description); err != nil {
			return nil, err
		}

		// Add to exercises set
		exercises.Exercises = append(exercises.Exercises, exercise)

	}
	// Build the total count query.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count.
	var total int
	if err = db.db.QueryRow(stmtSelectCount).Scan(&total); err != nil {
		return nil, err
	}
	exercises.Total = total
	// return exercises and nil
	return exercises, nil

}

// GetByID retrieves a exercise by its id.
func (db *Database) GetByID(id int) (*Exercise, error) {

	// Creat a new Exercise.
	exercise := &Exercise{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&exercise.ID, &exercise.ExerciseName, &exercise.MuscleGroup, &exercise.Description)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrExerciseNotFound
	case err != nil:
		return nil, err
	}

	// return exercise and nil
	return exercise, nil

}

// GetByName retrieves a exercise by its exercise_name.
func (db *Database) GetByName(name string) (*Exercise, error) {

	// Create a new Exercise.
	exercise := &Exercise{}

	// Execute eh query
	err := db.db.QueryRow(stmtSelectByName, name).Scan(&exercise.ID, &exercise.ExerciseName, &exercise.MuscleGroup, &exercise.Description)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrExerciseNotFound
	case err != nil:
		return nil, err
	}
	// return exercise and nil.
	return exercise, nil
}

// GetByMuscleGroup retrieves a group of exercises by muscle_group
func (db *Database) GetByMuscleGroup(muscleGroup string) (*Exercises, error) {
	// Build the full query.
	query := fmt.Sprintf(stmtSelectByMuscleGroup)

	// Create a new Exercises.
	exercises := &Exercises{
		Exercises: []*Exercise{},
	}

	// Execute the querey.
	rows, err := db.db.Query(query, muscleGroup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through the exercise row.
	for rows.Next() {
		// Create a new Exercise
		exercise := &Exercise{}

		// Scan row values into Exercise struct.
		if err := rows.Scan(&exercise.ID, &exercise.ExerciseName, &exercise.MuscleGroup, &exercise.Description); err != nil {
			return nil, err
		}

		// Add to exercises set
		exercises.Exercises = append(exercises.Exercises, exercise)

	}
	// Build the total count query.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count.
	var total int
	if err = db.db.QueryRow(stmtSelectCount).Scan(&total); err != nil {
		return nil, err
	}
	exercises.Total = total
	// return exercises and nil
	return exercises, nil
}
