package middleware

import (
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
)

var protectedRoutes []endpoints.Route

func AddProtectedRoute(route endpoints.Route) {
	protectedRoutes = append(protectedRoutes, route)
}

// AuthMiddleware is a middleware that checks if the user supplied a valid token
// the endpoints themselves are responsible for checking if the user has the correct role
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is to a protected route
		for _, route := range protectedRoutes {
			if c.Request.Method == route.Method && c.Request.URL.Path == route.Path {
				validateAuthHeader(c)
			}
		}

		c.Next()
	}
}

func validateAuthHeader(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	expectedPrefix := "Bearer "

	validBearer := len(c.GetHeader("Authorization")) > len(expectedPrefix) && c.GetHeader("Authorization")[:len(expectedPrefix)] == expectedPrefix

	if !validBearer {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	token := c.GetHeader("Authorization")[len(expectedPrefix):]

	if token == "" {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	_, err := models.ParseToken(token)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}
}
