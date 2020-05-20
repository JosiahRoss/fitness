package auth

import (
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
