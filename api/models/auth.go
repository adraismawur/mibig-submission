package models

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
)

type Token struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
	jwt.RegisteredClaims
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
