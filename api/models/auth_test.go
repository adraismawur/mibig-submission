package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJlbWFpbCI6InRlc3RAbG9jYWxob3N0IiwicGFzc3dvcmQiOiIiLCJhY3RpdmUiOmZhbHNlLCJyb2xlcyI6W3sicm9sZSI6ImFkbWluIn0seyJyb2xlIjoicmV2aWV3ZXIifSx7InJvbGUiOiJzdWJtaXR0ZXIifV0sImluZm8iOnsidXNlcl9pZCI6MCwiYWxpYXMiOiIiLCJuYW1lIjoiIiwiY2FsbF9uYW1lIjoiIiwib3JnYW5pemF0aW9uMSI6IiIsIm9yZ2FuaXphdGlvbjIiOiIiLCJvcmdhbml6YXRpb24zIjoiIiwib3JjX2lkIjoiIiwiUHVibGljIjpmYWxzZX19LCJpc3MiOiJtaWJpZy1zdWJtaXNzaW9uLWJlIiwic3ViIjoidGVzdEBsb2NhbGhvc3QiLCJhdWQiOlsibWliaWctc3VibWlzc2lvbi1mZSJdLCJleHAiOjE3NTI4NjMyMzIwLCJuYmYiOjE3NDcyMzU3NDQsImlhdCI6MTc0NzIzNTc0NH0.6tmrr83TXkkfcrlybBH8TzvtFZCEsNXoMWaXePqIBe50Z3eX1Zm7UwuSnglEJvYvBnEynJ-74eojgeqrX1qoGw"

	parsedToken, err := ParseToken(token)

	assert.Nil(t, err)

	expectedToken := Token{
		User: User{
			ID:     0,
			Email:  "test@localhost",
			Active: false,
			Roles: []UserRole{
				{
					Role: Admin,
				},
				{
					Role: Reviewer,
				},
				{
					Role: Submitter,
				},
			},
			Info: UserInfo{
				Public: false,
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
	}

	assert.Equal(t, expectedToken, parsedToken)
}
