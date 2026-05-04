package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MakeJWT creates a new signed JSON Web Token
func MakeJWT(userID string, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)

	// Create the standard claims payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "symbiosisos",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID, // We store the Postgres UUID as the subject
	})

	// Sign the token with our secret key
	return token.SignedString(signingKey)
}

// ValidateJWT parses the token string, verifies the signature, and returns the User ID (Subject)
func ValidateJWT(tokenString string, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}

	// Parse the token and verify the signing method
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("unauthorized")
	}

	// Extract the user ID we saved as the Subject earlier
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("invalid token subject")
	}

	return userIDString, nil
}

// GetBearerToken extracts the raw JWT string from the HTTP Authorization header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header included")
	}

	// The header should look like: "Bearer eyJhbG..."
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
