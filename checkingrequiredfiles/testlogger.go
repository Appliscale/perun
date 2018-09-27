package checkingrequiredfiles

import (
	"strings"

	"fmt"
)

type TestLogger struct {
	Quiet              bool
	Yes                bool
	Verbosity          Verbosity
	resourceValidation []*ResourceValidation
}

type ResourceValidation struct {
	ResourceName string
	Errors       []string
}

type Verbosity int

const (
	TRACE Verbosity = iota
	DEBUG
	INFO
	ERROR
	WARNING
)

var verboseModes = [...]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"ERROR",
	"WARNING",
}

func (verbosity Verbosity) String() string {
	return verboseModes[verbosity]
}
func CreateDefaultLogger() TestLogger {
	return TestLogger{
		Quiet:     false,
		Yes:       false,
		Verbosity: INFO,
	}
}

func (logTest *TestLogger) Always(message string) {
	fmt.Println(message)
}

// Log error.
func (logTest *TestLogger) Warning(warning string) {
	logTest.log(WARNING, warning)
}

// Log error.
func (logTest *TestLogger) Error(err string) {
	logTest.log(ERROR, err)
}

// Log info.
func (logTest *TestLogger) Info(info string) {
	logTest.log(INFO, info)
}

// Log debug.
func (logTest *TestLogger) Debug(debug string) {
	logTest.log(DEBUG, debug)
}

// Log trace.
func (logTest *TestLogger) Trace(trace string) {
	logTest.log(TRACE, trace)
}

// Log validation error.
func (resourceValidation *ResourceValidation) AddValidationError(error string) {
	resourceValidation.Errors = append(resourceValidation.Errors, error)
}

func (logTest *TestLogger) log(verbosity Verbosity, message string) {
	if !logTest.Quiet && verbosity >= logTest.Verbosity {
		fmt.Println(verbosity.String() + ": " + message)
	}
}

// Print validation error.
func (logTest *TestLogger) PrintValidationErrors() {
	if !logTest.Quiet {
		for _, resourceValidation := range logTest.resourceValidation {
			if len(resourceValidation.Errors) != 0 {
				fmt.Println(resourceValidation.ResourceName)
				for _, err := range resourceValidation.Errors {
					fmt.Println("        ", err)
				}
			}
		}
	}
}

func (logTest *TestLogger) HasValidationErrors() bool {
	for _, resourceValidation := range logTest.resourceValidation {
		if len(resourceValidation.Errors) > 0 {
			return true
		}
	}
	return false
}

// AddResourceForValidation : Adds resource for validation
func (logTest *TestLogger) AddResourceForValidation(resourceName string) *ResourceValidation {
	resourceValidation := &ResourceValidation{
		ResourceName: resourceName,
	}
	logTest.resourceValidation = append(logTest.resourceValidation, resourceValidation)

	return resourceValidation
}

// Set logger verbosity.
func (logTest *TestLogger) SetVerbosity(verbosity string) {
	for index, element := range verboseModes {
		if strings.ToUpper(verbosity) == element {
			logTest.Verbosity = Verbosity(index)
		}
	}
}

func (logTest *TestLogger) GetInput(message string, v ...interface{}) error {
	fmt.Printf("%s: ", message)
	_, err := fmt.Scanln(v...)
	if err != nil {
		return err
	}
	return nil
}
