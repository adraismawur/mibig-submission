package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QiLCJyb2xlIjowLCJpc3MiOiJtaWJpZy1zdWJtaXNzaW9uLWJlIiwic3ViIjoidGVzdCIsImF1ZCI6WyJtaWJpZy1zdWJtaXNzaW9uLWZlIl0sImV4cCI6MTc1MjEzNzQ3NzMsIm5iZiI6MTc0MjU2NDM0NywiaWF0IjoxNzQyNTY0MzQ3fQ.M5rluPSqBrYz87uI-GJsyLWh0ufyIhBiTqrHlpjzG2TzgM-abYu6r1YsYkQvas6D-nMtqlnSM6exxzkk92T0yQ"

	parsedToken, _ := ParseToken(token)

	assert.Equal(t, Token{
		Email: "test",
		Role:  0,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "mibig-submission-be",
			Subject:   "test",
			Audience:  []string{"mibig-submission-fe"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(17521374773, 0)), // in the year 2525, if man is still alive, they may find... the expiration date of this token
			IssuedAt:  jwt.NewNumericDate(time.Unix(1742564347, 0)),
			NotBefore: jwt.NewNumericDate(time.Unix(1742564347, 0)),
		},
	}, parsedToken)
}
