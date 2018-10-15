// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logger provides logger tool for perun for control standard I/O
// usage.
package logger

import (
	"fmt"
	"strings"
)

// Logger contains information type of logger tool.
type LoggerInt interface {
	Always(message string)
	Warning(warning string)
	Error(err string)
	Info(info string)
	Debug(debug string)
	Trace(trace string)
	GetInput(message string, v ...interface{}) error
	PrintValidationErrors()
	HasValidationErrors() bool
	AddResourceForValidation(resourceName string) *ResourceValidation
	SetVerbosity(verbosity string)
}
type Logger struct {
	Quiet              bool
	Yes                bool
	Verbosity          Verbosity
	resourceValidation []*ResourceValidation
}

// ResourceValidation contains name of resource and errors.
type ResourceValidation struct {
	ResourceName string
	Errors       []string
}

// Verbosity - type of logger.
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

// Create default logger.
func CreateDefaultLogger() Logger {
	return Logger{
		Quiet:     false,
		Yes:       false,
		Verbosity: INFO,
	}
}

// Create quiet logger.
func CreateQuietLogger() Logger {
	return Logger{
		Quiet:     true,
		Yes:       false,
		Verbosity: INFO,
	}
}

// Log always - no matter the verbosity level.
func (logger *Logger) Always(message string) {
	fmt.Println(message)
}

// Log error.
func (logger *Logger) Warning(warning string) {
	logger.log(WARNING, warning)
}

// Log error.
func (logger *Logger) Error(err string) {
	logger.log(ERROR, err)
}

// Log info.
func (logger *Logger) Info(info string) {
	logger.log(INFO, info)
}

// Log debug.
func (logger *Logger) Debug(debug string) {
	logger.log(DEBUG, debug)
}

// Log trace.
func (logger *Logger) Trace(trace string) {
	logger.log(TRACE, trace)
}

// Log validation error.
func (resourceValidation *ResourceValidation) AddValidationError(error string) {
	resourceValidation.Errors = append(resourceValidation.Errors, error)
}

// Get input from command line.
func (logger *Logger) GetInput(message string, v ...interface{}) error {
	fmt.Printf("%s: ", message)
	_, err := fmt.Scanln(v...)
	if err != nil {
		return err
	}
	return nil
}
func (logger *Logger) log(verbosity Verbosity, message string) {
	if !logger.Quiet && verbosity >= logger.Verbosity {
		fmt.Println(verbosity.String() + ": " + message)
	}
}

// Print validation error.
func (logger *Logger) PrintValidationErrors() {
	if !logger.Quiet {
		for _, resourceValidation := range logger.resourceValidation {
			if len(resourceValidation.Errors) != 0 {
				fmt.Println(resourceValidation.ResourceName)
				for _, err := range resourceValidation.Errors {
					fmt.Println("        ", err)
				}
			}
		}
	}
}

// HasValidationErrors checks if resource has errors. It's used in validateResources().
func (logger *Logger) HasValidationErrors() bool {
	for _, resourceValidation := range logger.resourceValidation {
		if len(resourceValidation.Errors) > 0 {
			return true
		}
	}
	return false
}

// AddResourceForValidation : Adds resource for validation. It's used in validateResources().
func (logger *Logger) AddResourceForValidation(resourceName string) *ResourceValidation {
	resourceValidation := &ResourceValidation{
		ResourceName: resourceName,
	}
	logger.resourceValidation = append(logger.resourceValidation, resourceValidation)

	return resourceValidation
}

// Set logger verbosity.
func (logger *Logger) SetVerbosity(verbosity string) {
	for index, element := range verboseModes {
		if strings.ToUpper(verbosity) == element {
			logger.Verbosity = Verbosity(index)
		}
	}
}

// Check if verbosity is one of the given types.
func IsVerbosityValid(verbosity string) bool {
	switch verbosity {
	case
		"TRACE",
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR":
		return true
	}
	return false
}
