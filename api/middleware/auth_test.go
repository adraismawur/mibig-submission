package middleware

import (
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	AddProtectedRoute(endpoints.Route{
		Method:  http.MethodGet,
		Path:    "/test",
		Handler: nil,
	})

	c := util.CreateMockGinGetRequest("/test")

	middleware := AuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status(), "Status code should be 401")
}
