package db

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConnectDefaults(t *testing.T) {
	config.Init()

	Connect()

	assert.NotNil(t, DB)
}

func TestConnectPostgres(t *testing.T) {
	// set environment variables
	os.Setenv("DB_DIALECT", "postgres")
	os.Setenv("DB_USER", "postgres")

	config.Init()

	Connect()

	assert.NotNil(t, DB, "DB using postgres dialect should not be nil")
}

func TestConnectUnsupported(t *testing.T) {
	// set environment variables
	os.Setenv("DB_DIALECT", "unsupported")

	config.Init()

	assert.Panics(t, func() {
		Connect()
	}, "Connect should panic when using unsupported dialect")
}
