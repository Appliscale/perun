// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

// Create default logger.
func CreateDefaultLogger() Logger {
	return Logger{
		Quiet:     false,
		Yes:       false,
		Verbosity: INFO,
	}
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
func (logger *Logger) ValidationError(elementName string, error string) {
	logger.validationErrors = append(logger.validationErrors, "\""+elementName+"\" "+error)
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
		for _, err := range logger.validationErrors {
			fmt.Println(err)
		}
	}
}

// Set logger verbosity.
func (logger *Logger) SetVerbosity(verbosity string) {
	for index, element := range verboseModes {
		if strings.ToUpper(verbosity) == element {
			logger.Verbosity = Verbosity(index)
		}
	}
}
