package middleware

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthMiddlewareNoToken(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test")

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status(), "Status code should be 401")
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Email:            "",
		Role:             models.Admin,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	signedToken, _ := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	c.Request.Header.Add("Authorization", "Bearer "+signedToken)

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
}

func TestAuthMiddlewareWrongSecret(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Email:            "",
		Role:             models.Admin,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	signedToken, _ := token.SignedString([]byte("wrong secret"))

	c.Request.Header.Add("Authorization", "Bearer "+signedToken)

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status(), "Status code should be 401")
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test")

	c.Request.Header.Add("Authorization", "Bearer ")

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status(), "Status code should be 401")
}

func TestAuthMiddlewareWrongRole(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Email:            "",
		Role:             models.Submitter,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	signedToken, _ := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	c.Request.Header.Add("Authorization", "Bearer "+signedToken)

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}

func TestAuthMiddlewareParameterizedRoute(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test/:id", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test/1")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Email:            "",
		Role:             models.Admin,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	signedToken, _ := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	c.Request.Header.Add("Authorization", "Bearer "+signedToken)

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
}

func TestAuthMiddlewareParameterizedRouteWrongRole(t *testing.T) {
	AddProtectedRoute(http.MethodGet, "/test/:id", models.Admin)

	c, _ := util.CreateMockGinGetRequest("/test/1")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Token{
		Email:            "",
		Role:             models.Submitter,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	signedToken, _ := token.SignedString([]byte(config.Envs["JWT_SECRET"]))

	c.Request.Header.Add("Authorization", "Bearer "+signedToken)

	middleware := AuthMiddleware()
	middleware(c)

	//assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}
