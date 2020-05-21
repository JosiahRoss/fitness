package exercises

import (
	"fitness/api/errors"
	"fitness/api/middleware/auth"
	"net/http"

	apictx "fitness/api/context"
	servexercises "fitness/services/exercises"

	"github.com/beeker1121/httprouter"
)

// Exercise defines a exercise
type Exercise struct {
	ID           int    `json:"id"`
	ExerciseName string `json:"exercise_name"`
	MuscleGroup  string `json:"muscle_group"`
	Description  string `json:"description"`
}

// Meta defines teh response top level meta object.
type Meta struct {
	Total int `json:"total"`
}

// ResultGetAll defines teh response data for the HandlerGet handler.
type ResultGetAll struct {
	Data []*Exercise `json:"data"`
	Meta Meta        `json:"meta"`
}

// ResultGetExercise teh response data for the HandlerGetExercise.
type ResultGetExercise struct {
	Data *Exercise `json:"data"`
}

// New creates the routes for the exercises of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/vi/exercises", auth.AuthenticateEndpoint(ac, HandleGetAll(ac)))
}

// HandleGetAll handles the /api/v1/exercises GetAll route of the API.
func HandleGetAll(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get this user form the request context.
		user, err := auth.GetUserFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new GetParams.
		params := &servexercises.
	}
}
