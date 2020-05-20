package signup

import (
	"encoding/json"
	"net/http"

	"github.com/beeker1121/httprouter"

	apictx "fitness/api/context"
	"fitness/api/errors"
	"fitness/api/middleware/auth"
	"fitness/api/render"
	serverrors "fitness/services/errors"
	"fitness/services/users"
)

// ResultPost defines teh response data fro the  HandlePost handler.
type ResultPost struct {
	Data string `json:"data"`
}

// New creates the routes for the signup endpoints of  the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/signup", HandlePost(ac))

}

// HandlePost heandles the
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse the parameters from the request body.
		var params users.NewParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Create the user
		user, err := ac.Services.Users.New(&params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("users.New() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Issue a new JWT for this user.
		token, err := auth.JWT(ac, user.Password, user.ID)
		if err != nil {
			ac.Logger.Printf("auth.NewJWT() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: token,
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

	}

}
