package config

import (
	"log"
	"os"
)

type EnvVar struct {
	DefaultValue string
	Optional     bool
}

var Envs map[string]string

var defaults = map[string]string{
	"DB_DIALECT": "sqlite",
	"DB_PATH":    "/tmp/test.db",
	"DB_HOST":    "localhost",
	"DB_PORT":    "5432",
	"DB_DBNAME":  "mibig_submission",
	"DB_USER":    "postgresql",
	"DB_PASS":    "",
}
var envLog = log.New(os.Stdout, "[env] ", log.Lmsgprefix)

func Init() {
	envLog.Println("Environment variables:")

	Envs = make(map[string]string)
	// replace with actual values from system environment
	for name := range defaults {
		sysValue := os.Getenv(name)
		if sysValue != "" {
			Envs[name] = sysValue
		} else {
			Envs[name] = defaults[name]
		}
		envLog.Printf("%s: %s\n", name, Envs[name])
	}
}
