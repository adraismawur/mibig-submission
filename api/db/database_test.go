package db

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectDefaults(t *testing.T) {
	db, err := Connect()

	assert.NotNil(t, db)
	assert.Nil(t, err)
}

func TestConnectUnsupported(t *testing.T) {
	// set environment variables
	config.OverrideEnv(config.EnvDbDialect, "unsupported")

	assert.Panics(t, func() {
		Connect()
	}, "Connect should panic when using unsupported dialect")
}
