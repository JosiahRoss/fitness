package auth

import (
	"context"
	"fitness/api/errors"
	"fitness/database/users"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	apictx "gotodo/api/context"
)

// key is the key type used by this package for the request context.
type key int

// AuthKey is the key used for storing and retrieving the member data from the
// request context.
var AuthKey key = 1

// TokenClaims defines the custom claims we use for the JWT.
type TokenClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// NewJWT creates and returns a new signed JWT.
func NewJWT(ac *apictx.Context, userPassword string, mid int) (string, error) {
	// Set expiry time.
	issued := time.Now()
	expires := issued.Add(time.Minute * ac.Config.JWTExpiryTime)

	// Create the claims.
	claims := &TokenClaims{
		mid,
		jwt.StandardClaims{
			IssuedAt:  issued.Unix(),
			ExpiresAt: expires.Unix(),
		},
	}

	// Create and sign the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(GetJWTSigningKey(ac.Config.JWTSecret, userPassword))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// AuthenticateEndpoint is the middleware for authenticating API requests.
//
// This function will first try to determine the type of authorization being
// requested, and then either authorize via a JWT or an API key.
//
// JWTs are passed via the Authorization header as a Bearer token.
//
// API keys should be passed via the Authorization header using Basic Auth.
//
// Currently, the only supported authorization method is via JWTs.
func AuthenticateEndpoint(ac *apictx.Context, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &users.User{}
		var err error

		// Get the Authorization header.
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")

		// Check for either Bearer or Basic authoriation type.
		if authHeader[0] != "Bearer" && authHeader[0] != "Basic" {
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", ErrUnauthorized.Error()))
			return
		}

		if len(authHeader) == 2 && authHeader[0] == "Bearer" {
			// Try authorization via JWT Authorization Bearer header first.
			user, err = GetUserFromJWT(ac, authHeader[1])
			if err == ErrJWTUnauthorized {
				ac.Logger.Println("API authorization via JWT failure")
				errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", err.Error()))
				return
			} else if err != nil {
				ac.Logger.Printf("auth.GetUserFromJWT() error: %s\n", err)
				errors.Default(ac.Logger, w, errors.ErrInternalServerError)
				return
			}
		} else {
			// Get the user from the API key.
			ac.Logger.Println("API key authorization not implemented")
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", "API key authorization not implemented"))
			return
		}

		// Pass member to request context and call next handler.
		ctx := context.WithValue(r.Context(), AuthKey, user)
		h(w, r.WithContext(ctx))
	}
}
