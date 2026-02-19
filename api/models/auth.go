package models

import (
	"errors"
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Token struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

func GetUserFromContext(c *gin.Context) (*User, error) {
	bearerToken, err := GetAuthHeaderToken(c)

	if err != nil {
		slog.Error("Could not get bearer token from request")
		return nil, err
	}

	parsedToken, err := ParseToken(bearerToken)

	if err != nil {
		slog.Error("Could not parse bearer token")
		return nil, err
	}

	return &parsedToken.User, nil
}

func GetAuthHeaderToken(c *gin.Context) (string, error) {
	expectedPrefix := "Bearer "

	// check there is a token behind prefix and check if the prefix is correct
	validBearer := len(c.GetHeader("Authorization")) > len(expectedPrefix) && c.GetHeader("Authorization")[:len(expectedPrefix)] == expectedPrefix

	if !validBearer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return "", errors.New("invalid bearer prefix")
	}

	// get the actual token
	bearerToken := c.GetHeader("Authorization")[len(expectedPrefix):]

	if bearerToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", errors.New("empty bearer token")
	}

	return bearerToken, nil
}

func ParseToken(token string) (Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		jwtSecret, err := config.GetConfig(config.EnvJwtSecret)

		if err != nil {
			return nil, err
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		slog.Error(fmt.Sprintf("[auth] Error parsing token: %s", err))
		return Token{}, err
	}

	if claims, ok := parsedToken.Claims.(*Token); ok && parsedToken.Valid {
		return *claims, nil
	}

	slog.Error(fmt.Sprintf("[auth] Error parsing token: Could not parse claims type"))
	panic("Invalid token")
}

func GetUserFromToken(token string) (*User, error) {
	parsedToken, err := ParseToken(token)

	if err != nil {
		return nil, err
	}

	return &parsedToken.User, nil
}

func GenerateToken(user User) (string, error) {
	issuedAt := time.Now()

	jwtSecret, err := config.GetConfig(config.EnvJwtSecret)

	if err != nil {
		return "", err
	}

	jwtLifetime, err := config.GetConfig(config.EnvJwtLifetime)

	if err != nil {
		return "", err
	}

	lifetime, err := strconv.ParseInt(jwtLifetime, 10, 64)

	if err != nil {
		slog.Error(fmt.Sprintf("[token] Error parsing JWT lifetime '%s'", jwtLifetime))
		return "", err
	}

	expirationTime := issuedAt.Add(time.Duration(lifetime) * time.Second)

	claims := Token{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    "mibig-entry-be",
			Subject:   user.Email,
			Audience:  []string{"mibig-entry-fe"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		slog.Error(fmt.Sprintf("[token] Error signing token: %s", err))
		return "", err
	}

	return tokenString, err
}
