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

var Envs map[string]string

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

func Init() {
	slog.Info("[env] Environment variables:")

	Envs = make(map[string]string)
	// replace with actual values from system environment
	for name := range defaults {
		sysValue := os.Getenv(name)
		if sysValue != "" {
			Envs[name] = sysValue
		} else {
			Envs[name] = defaults[name]
		}
		slog.Info(fmt.Sprintf("%s: %s", name, Envs[name]))
	}
}
