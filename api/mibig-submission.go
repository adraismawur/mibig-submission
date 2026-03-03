// Package main contains only the main function, which is the entry point of the application.
package main

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/functions"
	"log/slog"
	"os"
)

type RunFunction string

const (
	Run          RunFunction = "run"
	CheckExports RunFunction = "check-exports"
)

var functionDescriptions = map[RunFunction]string{
	Run:          "Run MIBiG Submission Portal API",
	CheckExports: "Check MIBiG entry exports",
}

type Args struct {
	Function RunFunction
}

// main is the entry point of the application
func main() {
	if len(os.Args) == 1 {
		slog.Error("[main] Missing required argument")
		slog.Error("[main]\tUsage: mibig-submission function")
		slog.Error("[main]\tAvailable functions:")
		for key, value := range functionDescriptions {
			slog.Error(fmt.Sprintf("[main]\t\t%s: %s", key, value))
		}
		return
	}

	os.Args = os.Args[1:]

	args := &Args{
		Function: RunFunction(os.Args[0]),
	}

	slog.Info(fmt.Sprintf("[main] Starting main function: %v", args.Function))

	switch args.Function {
	case Run:
		functions.Run()
		break
	case CheckExports:
		functions.CheckExports()
		break
	}
}
