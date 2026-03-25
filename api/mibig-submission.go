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
	StartApi             RunFunction = "start-api"
	CheckExports         RunFunction = "check-exports"
	StartAntismashWorker             = "antismash-worker"
)

var functionDescriptions = map[RunFunction]string{
	StartApi:             "StartApi MIBiG Submission Portal API",
	CheckExports:         "Check MIBiG entry exports",
	StartAntismashWorker: "Start antismash worker",
}

type Args struct {
	Function RunFunction
}

func printFunctionError() {
	slog.Error("[main]\tAvailable functions:")
	for key, value := range functionDescriptions {
		slog.Error(fmt.Sprintf("[main]\t\t%s: %s", key, value))
	}
}

// main is the entry point of the application
func main() {
	if len(os.Args) == 1 {
		slog.Error("[main] Missing required argument")
		slog.Error("[main]\tUsage: mibig-submission function")
		printFunctionError()
		return
	}

	os.Args = os.Args[1:]

	args := &Args{
		Function: RunFunction(os.Args[0]),
	}

	slog.Info(fmt.Sprintf("[main] Starting main function: %v", args.Function))

	switch args.Function {
	case StartApi:
		functions.StartApi()
		break
	case CheckExports:
		functions.CheckExports()
		break
	case StartAntismashWorker:
		functions.StartAntismashWorker()
		break
	default:
		slog.Error(fmt.Sprintf("[main] Unknown function: %v", args.Function))
		printFunctionError()
	}
}
