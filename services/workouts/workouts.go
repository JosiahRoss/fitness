package workouts

import (
	"fitness/database"
	dbworkouts "fitness/database/workouts"
	"fitness/services/errors"
)

// Service defines the workouts service.
type Service struct {
	db *database.Database
}

// New returns a new workouts service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Workout defines a workout.
type Workout dbworkouts.Workout

// Workouts defines a set of workouts.
type Workouts struct {
	Workouts []*Workout `json:"workouts"`
	Total    int        `json:"total"`
}

// NewParams defines the parameters for the New method.
type NewParams dbworkouts.NewParams

// New creates a new workout.
func (s *Service) New(sid int, params *NewParams) (*Workout, error) {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check

}
