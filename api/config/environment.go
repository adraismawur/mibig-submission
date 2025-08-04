// Package config provides ways to read the runtime configuration of the
// application using environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"
)

type EnvVar struct {
	DefaultValue string
	Optional     bool
}

// Envs is a map of environment variables and their values
// by default, it is initialized with the default values
var Envs = map[string]string{
	"SERVER_PORT":  ":8000",
	"DB_DIALECT":   "sqlite",
	"DB_PATH":      "/tmp/test.db",
	"DB_HOST":      "localhost",
	"DB_PORT":      "5432",
	"DB_DBNAME":    "mibig_submission",
	"DB_USER":      "postgres",
	"DB_PASS":      "",
	"JWT_SECRET":   "CHANGEME",
	"JWT_LIFETIME": "86400",
	"DATA_PATH":    "data",
}

// init initializes the Envs map with the environment variables,
// replacing any default values with actual values from the system environment.
func init() {
	slog.Info("[env] Environment variables:")

	// replace with actual values from system environment
	for name := range Envs {
		sysValue := os.Getenv(name)
		// if the environment variable is not set, use default value
		if sysValue == "" {
			slog.Info(fmt.Sprintf("%s: %s (DEFAULT)", name, Envs[name]))
			continue
		}

		Envs[name] = sysValue

		slog.Info(fmt.Sprintf("%s: %s", name, Envs[name]))
	}
}

// TODO: Add a function to validate the environment variables
// TODO: Add a bool/check to see if the environment variables are initialized
// necessary to prevent the application from running if the environment variables are not set
