package exercises

import (
	"fmt"
	"net/http"

	apictx "fitness/api/context"
	"fitness/api/errors"
	"fitness/api/render"

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
	router.GET("/api/v1/exercises", HandleGetAll(ac))

}

// HandleGetAll handles the /api/v1/exercises GetAll route of the API.
func HandleGetAll(ac *apictx.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get this user form the request context.
		// user, err := auth.GetUserFromRequest(r)
		// if err != nil {
		// 	errors.Default(ac.Logger, w, errors.ErrInternalServerError)
		// 	return
		// }

		// Try to get the Exercises.
		exercises, err := ac.Services.Exercises.GetAll()
		if err != nil {
			fmt.Println(err)
			ac.Logger.Printf("exercises.GetAll() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGetAll{
			Data: []*Exercise{},
			Meta: Meta{
				Total: exercises.Total,
			},
		}

		// Loop through the exercises.
		for _, t := range exercises.Exercises {
			// Copy the exercise Type over.
			exercise := &Exercise{
				ID:           t.ID,
				ExerciseName: t.ExerciseName,
				MuscleGroup:  t.MuscleGroup,
				Description:  t.Description,
			}
			result.Data = append(result.Data, exercise)
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {

			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

	}
}
