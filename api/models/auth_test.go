package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJlbWFpbCI6InRlc3RAbG9jYWxob3N0IiwicGFzc3dvcmQiOiIiLCJhY3RpdmUiOmZhbHNlLCJyb2xlcyI6W3sicm9sZSI6MH1dLCJpbmZvIjp7InVzZXJfaWQiOjAsImFsaWFzIjoiIiwibmFtZSI6IiIsImNhbGxfbmFtZSI6IiIsIm9yZ2FuaXphdGlvbjEiOiIiLCJvcmdhbml6YXRpb24yIjoiIiwib3JnYW5pemF0aW9uMyI6IiIsIm9yY19pZCI6IiIsIlB1YmxpYyI6ZmFsc2V9fSwiaXNzIjoibWliaWctc3VibWlzc2lvbi1iZSIsInN1YiI6InRlc3RAbG9jYWxob3N0IiwiYXVkIjpbIm1pYmlnLXN1Ym1pc3Npb24tZmUiXSwiZXhwIjoxNzUyODYzMjMyMCwibmJmIjoxNzQ3MjM1NzQ0LCJpYXQiOjE3NDcyMzU3NDR9.h4oCp9s1DlwgrQn4JoU8hj2vQbPOu6G8hhRoTuYMq_aXRQfIvQgdwXFwQNY6fHnEU9i0lOrdlRGC5oIxQazSzg"

	parsedToken, _ := ParseToken(token)

	assert.Equal(t, Token{
		User: User{
			ID:    1,
			Email: "test@localhost",
			Roles: []UserRole{
				{
					Role: Submitter,
				},
			},
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "mibig-submission-be",
			Subject:   "test@localhost",
			Audience:  []string{"mibig-submission-fe"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(17528632320, 0)), // in the year 2525, if man is still alive, they may find... the expiration date of this token
			IssuedAt:  jwt.NewNumericDate(time.Unix(1747235744, 0)),
			NotBefore: jwt.NewNumericDate(time.Unix(1747235744, 0)),
		},
	}, parsedToken)
}
