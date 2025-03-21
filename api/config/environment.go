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
var Envs map[string]string

// defaults is a map of environment variables and their default values
// these are the expected environment variables that the application will use
var defaults = map[string]string{
	"DB_DIALECT":   "sqlite",
	"DB_PATH":      "/tmp/test.db",
	"DB_HOST":      "localhost",
	"DB_PORT":      "5432",
	"DB_DBNAME":    "mibig_submission",
	"DB_USER":      "postgres",
	"DB_PASS":      "",
	"JWT_SECRET":   "CHANGEME",
	"JWT_LIFETIME": "86400",
}

// Init initializes the Envs map with the environment variables,
// replacing any default values with actual values from the system environment.
func Init() {
	slog.Info("[env] Environment variables:")

	Envs = make(map[string]string)
	// replace with actual values from system environment
	for name := range defaults {
		sysValue := os.Getenv(name)
		if sysValue != "" {
			Envs[name] = sysValue
		} else {
			// use default value if not set
			Envs[name] = defaults[name]
		}
		slog.Info(fmt.Sprintf("%s: %s", name, Envs[name]))
	}
}

// TODO: Add a function to validate the environment variables
// TODO: Add a bool/check to see if the environment variables are initialized
// necessary to prevent the application from running if the environment variables are not set
