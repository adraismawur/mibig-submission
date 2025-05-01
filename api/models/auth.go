package models

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"strconv"
	"time"
)

type Token struct {
	ID    uint   `json:"id"`
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

func GenerateToken(id uint, email string, role Role) (string, error) {
	issuedAt := time.Now()
	lifetime, err := strconv.ParseInt(config.Envs["JWT_LIFETIME"], 10, 64)
	expirationTime := issuedAt.Add(time.Duration(lifetime) * time.Second)

	if err != nil {
		slog.Error(fmt.Sprintf("[token] Error parsing JWT lifetime '%s'", config.Envs["JWT_LIFETIME"]))
		return "", err
	}

	claims := Token{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    "mibig-submission-be",
			Subject:   email,
			Audience:  []string{"mibig-submission-fe"},
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
