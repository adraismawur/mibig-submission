package config

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestInit(t *testing.T) {
	// tests run with defaults
	assert.Equal(t, Envs["DB_HOST"], "localhost")
	assert.Equal(t, Envs["DB_DIALECT"], "sqlite")
}
