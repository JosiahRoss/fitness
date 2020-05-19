package workouts

import (
	"database/sql"
	"fmt"
)

// Database defines the workouts database
type Database struct {
	db *sql.DB
}

// New creates a new workouts database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Workout defines a workout.
type Workout struct {
	ID         int `json:"id"`
	SessionID  int `json:"session_id"`
	ExerciseID int `json:"exersizer_id"`
	Weight     int `json:"weight"`
	Reps       int `json:"reps"`
	Sets       int `json:"sets"`
}

// Workouts defines a set of worokouts.
type Workouts struct {
	Workouts []*Workout `json:"workouts"`
	Total    int        `json:"total"`
}

const (
	// stmtInsert defines teh SQL statement to
	// insert a new Workoutsint the database
	stmtInsert = `
INSERT INTO workouts (session_id, exercise_id ,weight,reps,sets)
VALUES(?,?,?,?,?)
`

	// stmtSelect defines the SQL statement to
	// select a set of exersizes in a workout.
	stmtSelect = `
SELECT  id, session_id, exercise_id ,weight,reps,sets
FROM workouts
%s
LIMIT %v, %v
`

	// stmtSelectCount defines the SQL statement to
	// select the total number of workouts found for a
	// a given member, according to the filters.
	stmtSelectCount = `
SELECT COUNT(*)
FROM workouts
%s
`

	// stmtSelectByID defines the SQL statement to
	// select a Workout by its ID.
	stmtSelectByID = `
SELECT id, session_id, exercise_id ,weight,reps,sets
FROM workouts
WHERE id=?
`

	// stmtSelectBySessionID defines the SQL statment to
	// select a set of Workouts by its session_id
	stmtSelectBySessionID = `
SELECT id, session_id, exercise_id ,weight,reps,sets
FROM workouts
WHERE session_id=?
`

	// stmtSelectByExerciseID defines the SQL statment to
	// select a set of Workouts by its exercise_id.
	stmtSelectByExerciseID = `
SELECT id, session_id, exercise_id ,weight,reps,sets
FROM workouts
WHERE exercise_id=?
`

	// stmtUpdate defines the SQL statment to
	// update a workout.
	stmtUpdate = `
UPDATE workouts
SET %s
WHERE id=?
`
)

// NewParams defines the parameters for the new method.
type NewParams struct {
	ExerciseID int `json:"exersize_id"`
	Weight     int `json:"weight"`
	Reps       int `json:"reps"`
	Sets       int `json:"sets"`
}

// New creates anad returns a  new workout
func (db *Database) New(sid int, params *NewParams) (*Workout, error) {
	// Create a new Workout
	workout := &Workout{
		SessionID:  sid,
		ExerciseID: params.ExerciseID,
		Weight:     params.Weight,
		Reps:       params.Reps,
		Sets:       params.Sets,
	}

	// Create variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, workout.SessionID, workout.ExerciseID, workout.Weight, workout.Reps, workout.Sets); err != nil {
		return nil, err
	}

	// Get last insert ID.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// Define workout.ID
	workout.ID = int(id)

	return workout, nil
}

// GetParams defines the Get parameters.
type GetParams struct {
	ID         *int `json:"id"`
	SessionID  *int `json:"session_id"`
	ExerciseID *int `json:"exersizer_id"`
	Weight     int  `json:"weight"`
	Reps       int  `json:"reps"`
	Sets       int  `json:"sets"`
	Offset     int  `json:"offset"`
	Limit      int  `json:"limit"`
}

// Get a group of workouts
func (db *Database) Get(params *GetParams) (*Workouts, error) {
	// Create variables to hold the query fields
	// being filtered on their values.
	var queryFields string
	var queryValues []interface{}

	// Handle ID field.
	if params.ID != nil {
		if queryFields == "" {
			queryFields = "WHERE id=?"
		} else {
			queryFields += " AND id=?"
		}

		queryValues = append(queryValues, *params.ID)
	}
	// Handle SessionID field.
	if params.SessionID != nil {
		if queryFields == "" {
			queryFields = "WHERE session_id=?"
		} else {
			queryFields += " AND session_id=?"
		}
		queryValues = append(queryValues, *params.SessionID)
	}
	// Handle ExerciseID field.
	if params.ExerciseID != nil {
		if queryFields == "" {
			queryFields = "WHERE exercise_id=?"
		} else {
			queryFields += " AND exercise_id=?"
		}
		queryValues = append(queryValues, *params.ExerciseID)
	}
	// Build the full query.
	query := fmt.Sprintf(stmtSelect, queryFields, params.Offset, params.Limit)
	// Create a new Workouts.
	workouts := &Workouts{
		Workouts: []*Workout{},
	}

	// Execute the query.
	rows, err := db.db.Query(query, queryValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the workout rows.
	for rows.Next() {
		// Create a new Workout.
		workout := &Workout{}

		// Scan ros values into workout struct.
		if err := rows.Scan(&workout.ID, &workout.SessionID, &workout.ExerciseID, &workout.Weight, &workout.Reps, &workout.Sets); err != nil {
			return nil, err
		}

		// Add to workouts set
		workouts.Workouts = append(workouts.Workouts, workout)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, queryFields)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount, queryValues...).Scan(&total); err != nil {
		return nil, err
	}

	workouts.Total = total

	return workouts, nil

}

// GetByID retrieves a workout by its ID
func (db *Database) GetByID(id int) (*Workout, error) {
	// Create a new Workout.
	workout := &Workout{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&workout.SessionID, &workout.ExerciseID, &workout.Weight, &workout.Reps, &workout.Sets)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrLogNotFound
	case err != nil:
		return nil, err

	}
	return workout, nil

}

// GetBySessionID  gets and returns an set of workouts by
// session_id
func (db *Database) GetBySessionID(sid int) (*Workouts, error) {
	// Create a new slice of workouts
	workouts := &Workouts{
		Workouts: []*Workout{},
	}

	// Execute the query.
	rows, err := db.db.Query(stmtSelectBySessionID, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the workout rows.
	for rows.Next() {
		// Create a new Workout.
		workout := &Workout{}

		// Scan ros values into workout struct.
		if err := rows.Scan(&workout.ID, &workout.SessionID, &workout.ExerciseID, &workout.Weight, &workout.Reps, &workout.Sets); err != nil {
			return nil, err
		}

		// Add to workouts set
		workouts.Workouts = append(workouts.Workouts, workout)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	field := "WHERE session_id=?"
	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, field)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, err
	}

	workouts.Total = total

	return workouts, nil
}

// GetByExerciseID  gets and returns an set of workouts by
// exercise_id
func (db *Database) GetByExerciseID(eid int) (*Workouts, error) {
	// Create a new slice of workouts
	workouts := &Workouts{
		Workouts: []*Workout{},
	}

	// Execute the query.
	rows, err := db.db.Query(stmtSelectByExerciseID, eid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop though the workout rows.
	for rows.Next() {
		// Create a new Workout.
		workout := &Workout{}

		// Scan ros values into workout struct.
		if err := rows.Scan(&workout.ID, &workout.SessionID, &workout.ExerciseID, &workout.Weight, &workout.Reps, &workout.Sets); err != nil {
			return nil, err
		}

		// Add to workouts set
		workouts.Workouts = append(workouts.Workouts, workout)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	field := "WHERE exercise_id=?"
	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, field)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, err
	}

	workouts.Total = total

	return workouts, nil
}

// UpdateParams defines the parameters for the Update mehtod.
type UpdateParams struct {
	ExerciseID *int `json:"exercise_id"`
	Weight     *int `json:"weight"`
	Reps       *int `json:"reps"`
	Sets       *int `json:"sets"`
}

// Update updates a workout.
func (db *Database) Update(id int, params *UpdateParams) (*Workout, error) {
	// Create variables to hold the query fields
	// being updated and their new values.
	var queryFields string
	var queryValues []interface{}

	// Handle exercise_id field.
	if params.ExerciseID != nil {
		if queryFields == "" {
			queryFields = "exercise_id=?"
		} else {
			queryFields += ", exercise_id=?"
		}

		queryValues = append(queryValues, *params.ExerciseID)
	}

	// Handle weight field.
	if params.Weight != nil {
		if queryFields == "" {
			queryFields = "weight=?"
		} else {
			queryFields += ", weight=?"
		}

		queryValues = append(queryValues, *params.Weight)
	}

	// Handle reps field.
	if params.Reps != nil {
		if queryFields == "" {
			queryFields = "reps=?"
		} else {
			queryFields += ", reps=?"
		}

		queryValues = append(queryValues, *params.Reps)
	}

	// Handle sets field.
	if params.Sets != nil {
		if queryFields == "" {
			queryFields = "sets=?"
		} else {
			queryFields += ", sets=?"
		}

		queryValues = append(queryValues, *params.Sets)
	}

	// Check if the query is empty.
	if queryFields == "" {
		return db.GetByID(id)
	}

	// Build the full query.
	query := fmt.Sprintf(stmtUpdate, queryFields)
	queryValues = append(queryValues, id)

	// Execute the query.
	_, err := db.db.Exec(query, queryValues...)
	if err != nil {
		return nil, err
	}

	// Since the GetByID method is straight forward,
	// we can use this method to retrieve the updated
	// Workout. Anything more complicated should use the
	// original statement constants.
	return db.GetByID(id)
}
