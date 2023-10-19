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

// Package specification privides tools for downloading and parsing AWS
// CloudFormation Resource Specification.
package specification

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/Appliscale/perun/context"
)

// Specification contains information about specification - type and version.
type Specification struct {
	PropertyTypes                map[string]PropertyType
	ResourceSpecificationVersion string
	ResourceTypes                map[string]Resource
}

// PropertyType contains Documentation and map of Properties.
type PropertyType struct {
	Documentation string
	Properties    map[string]Property
}

// Property describes what should be in structure.
type Property struct {
	Documentation     string
	DuplicatesAllowed bool
	ItemType          string
	PrimitiveItemType string
	PrimitiveType     string
	Required          bool
	Type              string
	UpdateType        string
}

// Resource in specification contains documentation, map of attributes and properties.
type Resource struct {
	Documentation string
	Attributes    map[string]Attribute
	Properties    map[string]Property
}

// Attribute of item.
type Attribute struct {
	ItemType          string
	PrimitiveItemType string
	PrimitiveType     string
	Type              string
}

// IsSubproperty : Checks if it is subproperty.
func (property *Property) IsSubproperty() bool {
	if property.Type != "List" && property.Type != "Map" && len(property.Type) > 0 {
		return true
	}
	return false
}

// Download specification for region specified in config.
func GetSpecification(context *context.Context) (specification Specification, err error) {
	filePath, err := downloadSpecification(context)
	if err != nil {
		return specification, err
	}

	return GetSpecificationFromFile(filePath)
}

// Get specification from file. It's used in GetSpecification().
func GetSpecificationFromFile(specificationFilePath string) (specification Specification, err error) {
	specificationFile, err := ioutil.ReadFile(specificationFilePath)
	if err != nil {
		return specification, err
	}

	return parseSpecificationFile(specificationFile)
}

func downloadSpecification(context *context.Context) (filePath string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	specificationDir := user.HomeDir + "/.config/perun/specification"
	specificationFileUrl, err := context.Config.GetSpecificationFileURLForCurrentRegion()
	if err != nil {
		return
	}
	fileName := strings.Replace(specificationFileUrl, "https://", "", -1)
	fileName = strings.Replace(fileName, ".cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json", "", -1)
	specificationFilePath := specificationDir + "/" + fileName + ".json"

	if _, err := os.Stat(specificationFilePath); err == nil {
		return specificationFilePath, nil
	}

	if _, err := os.Stat(specificationDir); os.IsNotExist(err) {
		os.MkdirAll(specificationDir, os.ModePerm)
	}
	out, err := os.Create(specificationFilePath)
	if err != nil {
		return
	}
	defer out.Close()

	resp, err := http.Get(specificationFileUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return specificationFilePath, nil
}

func parseSpecificationFile(specificationFile []byte) (specification Specification, err error) {
	err = json.Unmarshal(specificationFile, &specification)
	if err != nil {
		return specification, err
	}

	return specification, nil
}
