package login

import (
	"encoding/json"
	apictx "fitness/api/context"
	"fitness/api/errors"
	"fitness/api/middleware/auth"
	"fitness/api/render"
	serverrors "fitness/services/errors"
	"fitness/services/users"
	"fmt"
	"net/http"

	"github.com/beeker1121/httprouter"
)

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data string `json:"data"`
}

// New creates the routes for the login enpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/login", HandlePost(ac))
}

// HandlePost handles the /api/v1/login POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var params users.LoginParams
		fmt.Println("trying to decode")
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			fmt.Println(err)
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		fmt.Println("decoded")

		// Try to log this user in.
		user, err := ac.Services.Users.Login(&params)
		fmt.Println(err)

		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err == users.ErrInvalidLogin {
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Printf("users.Login() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Issue a new JWT for this user.
		token, err := auth.NewJWT(ac, user.Password, user.ID)
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
