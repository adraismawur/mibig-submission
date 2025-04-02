package middleware

import (
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// routeMap is a map using the route path and method as the key and a route as the value
// key is in the form "method:path"
var routeMap = make(map[string]models.Role)

func AddProtectedRoute(route endpoints.Route, role models.Role) {
	// add the route to the map
	key := route.Method + ":" + route.Path
	routeMap[key] = role
}

// AuthMiddleware is a middleware that checks if the user supplied a valid token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		key := method + ":" + path

		// check if route is in protected routes
		expectedRole, ok := routeMap[key]

		if !ok {
			c.Next()
			return
		}

		var token models.Token

		if !validateAuthHeader(c, &token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}

		if token.Role != expectedRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
			c.Abort()
		}

		c.Next()
	}
}

func validateAuthHeader(c *gin.Context, token *models.Token) bool {
	if c.GetHeader("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return false
	}

	expectedPrefix := "Bearer "

	validBearer := len(c.GetHeader("Authorization")) > len(expectedPrefix) && c.GetHeader("Authorization")[:len(expectedPrefix)] == expectedPrefix

	if !validBearer {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return false
	}

	bearerToken := c.GetHeader("Authorization")[len(expectedPrefix):]

	if bearerToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return false
	}

	parsedToken, err := models.ParseToken(bearerToken)

	token = &parsedToken

	if err != nil {
		return false
	}

	return true
}
