package exercises

import (
	"fitness/database"
	dbexercises "fitness/database/exercises"
)

// Service defines the exercises service.
type Service struct {
	db *database.Database
}

// New returns a new exercises service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}

}

// Exercise defines a exercise.
type Exercise dbexercises.Exercise

// Exercises defines a set of exercises.
type Exercises struct {
	Exercises []*Exercise `json:"exercises"`
	Total     int         `json:"total"`
}

// GetParams defines the parameters for the New method.
//type GetParams dbexercises.GetParams

// GetAll gets a set of all Exercises.
func (s *Service) GetAll() (*Exercises, error) {
	// Try to pull the exercises from the database.
	dbex, err := s.db.Exercises.GetAll()
	if err != nil {
		return nil, err
	}

	// Create a new Exersizes.
	exercises := &Exercises{
		Exercises: []*Exercise{},
		Total:     dbex.Total,
	}

	// Loop through the set of exercises.
	for _, t := range dbex.Exercises {
		// Create a new Exercise.
		exercise := &Exercise{
			ID:           t.ID,
			ExerciseName: t.ExerciseName,
			MuscleGroup:  t.MuscleGroup,
			Description:  t.Description,
		}

		// Add to exercises set.
		exercises.Exercises = append(exercises.Exercises, exercise)
	}
	return exercises, nil
}
