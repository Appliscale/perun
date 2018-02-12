// Copyright 2017 Appliscale
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

package context

import (
	"errors"
	"testing"

	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
)

func parseCliArgumentsValidStub(cliArguments cliparser.CliArguments) cliArgumentsParser {
	return func(args []string) (cliparser.CliArguments, error) {
		return cliArguments, nil
	}
}

func getConfigurationValidStub(config configuration.Configuration) configurationReader {
	return func(cliparser.CliArguments, *logger.Logger) (configuration.Configuration, error) {
		return config, nil
	}
}

func parseCliArgumentsErroneous(args []string) (cliparser.CliArguments, error) {
	return cliparser.CliArguments{}, errors.New("")
}

func getConfigurationErroneous(cliparser.CliArguments, *logger.Logger) (configuration.Configuration, error) {
	return configuration.Configuration{}, errors.New("")
}

func TestCheckContextBody(t *testing.T) {
	t.Run("CLI arguments returned from cliArgumentsParser are the same as the ones contained in context", func(t *testing.T) {
		cliArguments := cliparser.CliArguments{}
		config := configuration.Configuration{}
		cliArgParserStub := parseCliArgumentsValidStub(cliArguments)
		confReaderStub := getConfigurationValidStub(config)
		context, _ := GetContext(cliArgParserStub, confReaderStub)
		assert.Equal(t, cliArguments, context.CliArguments)
	})

	t.Run("Config returned from configurationReader is the same as the one contained in context", func(t *testing.T) {
		cliArguments := cliparser.CliArguments{}
		config := configuration.Configuration{}
		cliArgsParserStub := parseCliArgumentsValidStub(cliArguments)
		confReaderStub := getConfigurationValidStub(config)
		context, _ := GetContext(cliArgsParserStub, confReaderStub)
		assert.Equal(t, config, context.Config)
	})
}

func TestCheckErroneousDependencyReturn(t *testing.T) {
	t.Run("Should return an error if cliArgumentsParser returns the error", func(t *testing.T) {
		config := configuration.Configuration{}
		confReaderStub := getConfigurationValidStub(config)
		_, err := GetContext(parseCliArgumentsErroneous, confReaderStub)
		assert.NotNil(t, err)
	})

	t.Run("Should return an error if configurationReader returns the error", func(t *testing.T) {
		cliArguments := cliparser.CliArguments{}
		cliArgParserStub := parseCliArgumentsValidStub(cliArguments)
		_, err := GetContext(cliArgParserStub, getConfigurationErroneous)
		assert.NotNil(t, err)
	})
}
