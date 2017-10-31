package cflogger

import (
	"fmt"
	"strings"
)

type Logger struct {
	Quiet            bool
	Yes              bool
	Verbosity        Verbosity
	validationErrors []string
}

type Verbosity int

const (
	TRACE Verbosity = iota
	DEBUG
	INFO
	ERROR
)

var verboseModes = [...]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"ERROR",
}

func (verbosity Verbosity) String() string {
	return verboseModes[verbosity]
}

func CreateDefaultLogger() Logger {
	return Logger {
		Quiet: false,
		Yes: false,
		Verbosity: INFO,
	}
}

func (logger *Logger) Error(err string) {
	logger.log(ERROR, err)
}

func (logger *Logger) Info(info string) {
	logger.log(INFO, info)
}

func (logger *Logger) Debug(debug string) {
	logger.log(DEBUG, debug)
}

func (logger *Logger) Trace(trace string) {
	logger.log(TRACE, trace)
}

func (logger *Logger) ValidationError(elementName string, error string) {
	logger.validationErrors = append(logger.validationErrors, "\""+elementName+"\" "+error)
}

func (logger *Logger) log(verbosity Verbosity, message string) {
	if !logger.Quiet && verbosity >= logger.Verbosity {
		fmt.Println(verbosity.String() + ": " + message)
	}
}

func (logger *Logger) PrintValidationErrors() {
	if !logger.Quiet {
		for _, err := range logger.validationErrors {
			fmt.Println(err)
		}
	}
}

func (logger *Logger) SetVerbosity(verbosity string) {
	for index, element := range verboseModes {
		if strings.ToUpper(verbosity) == element {
			logger.Verbosity = Verbosity(index)
		}
	}
}