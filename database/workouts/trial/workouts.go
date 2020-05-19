package workouts

import (
	"database/sql"
	"fmt"
	"time"
)

//Database difines the Workouts database.
type Database struct {
	db *sql.DB
}

// New creates a new workout database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Workout defines a Workouts
type Workout struct {
	ID          int       `json:"log_id"`
	WorkoutID   int       `json:"workout_id"`
	ExerciseID  int       `json:"exerzise_id"`
	UserID      int       `json:"user_id"`
	WorkoutDate time.Time `json:"workout_date"`
	Weight      int       `json:"weight"`
	Reps        int       `json:"reps"`
	Sets        int       `json:"sets"`
}

// Workouts define a set  of logs
type Workouts struct {
	Workouts []*Workout `json:"workouts"`
	Total    int        `json:"total"`
}

const (
	// stmtInsert defines teh SQL statement to
	// insert a new Workoutsint the database
	stmtInsert = `
INSERT INTO workout_logs (workout_id, user_id,workout_date,weight,reps,sets)
VALUES(?,?,?,?,?,?,?)
`

	// stmtSelect defines the SQL statement to
	// select a set of exersizes in a workout.
	stmtSelect = `
SELECT workout_id, user_id, workout_date,weight,reps,sets
FROM todos
%s
LIMIT %v, %v
`

	// stmtSelectCount defines the SQL statement to
	// select the total number of workouts found for a
	// a given member, according to the filters.
	stmtSelectCount = `
SELECT COUNT(*)
FROM todos
%s
`

	// stmtSelectByID defines the SQL statement to
	// select a Workout by its ID.
	stmtSelectByID = `
SELECT id, workout_id, user_id, workout_date, weight, reps, sets
FROM todos
WHERE id=?
`

	// stmtSelectByID difines the SQL statement to
	// select  a workout by workoutID and user ID
	stmtSelectByID = `
SELECT id, member_id, created, detail, completed
FROM todos
WHERE workout_id=? AND user_id=?
`
)

// NewParams defines the parameters for the new method.
type NewParams struct {
	ExersizeID int `json:"exersize_id"`
	Weight     int `json:"weight"`
	Reps       int `json:"reps"`
	Sets       int `json:"sets"`
}

// New creates a Workoutsfunc (db *Database) New(lid int, woid int, uid int, params *NewParams) (*Workouts error){
func (db *Database) New(woid int, eid int, uid int, params *NewParams) (*Workout, error) {
	// Create a new  Workouts
	workout := &Workout{
		WorkoutID:   woid,
		ExerciseID:  eid,
		UserID:      uid,
		WorkoutDate: time.Now(),
		Weight:      params.Weight,
		Reps:        params.Reps,
		Sets:        params.Sets,
	}

	// Create variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, workout.WorkoutID, workout.ExerciseID, workout.UserID, workout.WorkoutDate, workout.Weight, workout.Reps, workout.Sets); err != nil {
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

// GetParams defines the getparameters.
type GetParams struct {
	ID          *int       `json:"log_id"`
	WorkoutID   *int       `json:"workout_id"`
	ExerciseID  *int       `json:"exercise_id"`
	UserID      *int       `json:"user_id"`
	WorkoutDate *time.Time `json:"workout_date"`
	Offset      int        `json:"offset"`
	Limit       int        `json:"limit"`
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

	// Handle WorkoutID field.
	if params.WorkoutID != nil {
		if queryFields == "" {
			queryFields = "WHERE workout_id=?"
		} else {
			queryFields += " AND workout_id=?"
		}
		queryValues = append(queryValues, *params.WorkoutID)
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

	// Handle the UserId field.
	if params.UserID != nil {
		if queryFields == "" {
			queryFields = "WHERE user_id=?"
		} else {
			queryFields += " AND user_id=?"
		}
		queryValues = append(queryValues, *params.UserID)
	}

	// Handle the WorkoutDate field.
	if params.WorkoutDate != nil {
		if queryFields == "" {
			queryFields = "WHERE workout_date=?"
		} else {
			queryFields += " AND workout_date=?"

		}
		queryValues = append(queryValues, *params.WorkoutDate)
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
		if err := rows.Scan(&workout.ID, &workout.WorkoutID, &workout.ExerciseID, &workout.UserID, &workout.WorkoutDate, &workout.Weight, &workout.Reps, &workout.Sets); err != nil {
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
func (db *Database) GetByID(id, uid int) (*Workout, error) {
	// Create a new Workout.
	workout := &Workout{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&workout.WorkoutID, &workout.ExerciseID, &workout.UserID, &workout.WorkoutDate, &workout.Weight, &workout.Reps, &workout.Sets)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrLogNotFound
	case err != nil:
		return nil, err

	}
	return workout, nil

}

// GetByWorkoutIDUserID retrieves a workout by WorkoutID and UserID

func (db *Database) GetByWorkoutIDUserID(woid, uid int) *Workout {

}
