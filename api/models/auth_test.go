package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QiLCJyb2xlIjowLCJpc3MiOiJtaWJpZy1zdWJtaXNzaW9uLWJlIiwic3ViIjoidGVzdCIsImF1ZCI6WyJtaWJpZy1zdWJtaXNzaW9uLWZlIl0sImV4cCI6MTc0MjY1MDc0NywibmJmIjoxNzQyNTY0MzQ3LCJpYXQiOjE3NDI1NjQzNDd9.rFe76lSpH1oOQ6O7tTWjZo5kzQOJQP2v36hGDsx0o0RAvRlAJWRhEVA53rAgkfsONrv96XTY_nVxTQISjHW-vg"

	parsedToken := ParseToken(token)

	assert.Equal(t, Token{
		Email: "test",
		Role:  0,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "mibig-submission-be",
			Subject:   "test",
			Audience:  []string{"mibig-submission-fe"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(1742650747, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(1742564347, 0)),
			NotBefore: jwt.NewNumericDate(time.Unix(1742564347, 0)),
		},
	}, parsedToken)
}
