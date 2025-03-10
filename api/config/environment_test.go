package config

import (
	"github.com/go-playground/assert/v2"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// set environment variables

	os.Setenv("DB_HOST", "test")
	os.Setenv("DB_DIALECT", "postgres")

	Init()

	assert.Equal(t, Envs["DB_HOST"], "test")
	assert.Equal(t, Envs["DB_DIALECT"], "postgres")
}
