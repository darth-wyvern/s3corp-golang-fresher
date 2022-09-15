package jwt

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTInput represents a JWT input to generates a JWT token
type JWTInput struct {
	ID        int
	Email     string
	Role      string
	SecretKey string
	ExpiresIn time.Duration
}

// GenerateJWTToken generates a token with the given data and secret key
func GenerateJWTToken(input JWTInput) (jwt.Token, string, error) {
	if input.SecretKey == "" {
		return nil, "", fmt.Errorf("secret key cannot be empty")
	}
	if input.ID < 0 {
		return nil, "", fmt.Errorf("invalid id")
	}
	if input.ID < 0 {
		return nil, "", fmt.Errorf("invalid email")
	}

	claim := map[string]interface{}{
		"id":    input.ID,
		"email": input.Email,
		"role":  input.Role,
	}
	jwtauth.SetExpiryIn(claim, input.ExpiresIn)
	tokenAuth := jwtauth.New("HS256", []byte(input.SecretKey), nil)

	return tokenAuth.Encode(claim)
}
