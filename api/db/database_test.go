package db

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectDefaults(t *testing.T) {
	db := Connect()

	assert.NotNil(t, db)
}

func TestConnectPostgres(t *testing.T) {
	// set environment variables
	config.Envs["DB_DIALECT"] = "postgres"
	config.Envs["DB_USER"] = "postgres"

	db := Connect()

	assert.NotNil(t, db, "DB using postgres dialect should not be nil")
}

func TestConnectUnsupported(t *testing.T) {
	// set environment variables
	config.Envs["DB_DIALECT"] = "unsupported"

	assert.Panics(t, func() {
		Connect()
	}, "Connect should panic when using unsupported dialect")
}
