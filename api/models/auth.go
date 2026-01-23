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
		return []byte(config.Envs["JWT_SECRET"]), nil
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
	lifetime, err := strconv.ParseInt(config.Envs["JWT_LIFETIME"], 10, 64)
	expirationTime := issuedAt.Add(time.Duration(lifetime) * time.Second)

	if err != nil {
		slog.Error(fmt.Sprintf("[token] Error parsing JWT lifetime '%s'", config.Envs["JWT_LIFETIME"]))
		return "", err
	}

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

	tokenString, err := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	if err != nil {
		slog.Error(fmt.Sprintf("[token] Error signing token: %s", err))
		return "", err
	}

	return tokenString, err
}
