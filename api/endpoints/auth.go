package endpoints

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

// GetAuthEndpoint returns the auth endpoint, which is responsible for specifically handling authentication.
// This means acquiring a token (logging in) and refreshing a token.
func GetAuthEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "POST",
				Path:   "/login",
				Handler: func(c *gin.Context) {
					login(db, c)
				},
			},
		},
	}
}

func login(db *gorm.DB, c *gin.Context) {
	var loginRequest models.LoginRequest
	err := c.BindJSON(&loginRequest)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user := models.User{}

	tx := db.First(&user, "email = ?", loginRequest.Email)

	if tx.RowsAffected == 0 {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		slog.Info(fmt.Sprintf("User %s not found", loginRequest.Email))
		return
	}

	if !models.CheckPassword(loginRequest.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	issuedAt := time.Now()
	lifetime, err := strconv.ParseInt(config.Envs["JWT_LIFETIME"], 10, 64)
	expirationTime := issuedAt.Add(time.Duration(lifetime) * time.Second)

	if err != nil {
		slog.Error(fmt.Sprintf("[db] [env] Error parsing JWT lifetime '%s'", config.Envs["JWT_LIFETIME"]))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	claims := models.Token{
		user.Email,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    "mibig-submission-be",
			Subject:   user.Email,
			Audience:  []string{"mibig-submission-fe"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	if err != nil {
		slog.Error(fmt.Sprintf("Error signing token: %s", err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

// TODO: Implement refresh token
