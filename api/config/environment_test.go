package config

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestInit(t *testing.T) {
	// tests run with defaults
	assert.Equal(t, envs["DB_HOST"], "localhost")
	assert.Equal(t, envs["DB_DIALECT"], "sqlite")
}

func TestGetConfig(t *testing.T) {
	val, _ := GetConfig(EnvDbHost)
	assert.Equal(t, val, "localhost")
}
