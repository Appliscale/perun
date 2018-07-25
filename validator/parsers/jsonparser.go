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

package parsers

import (
	"github.com/Appliscale/jsonparser"
	"github.com/Appliscale/perun/validator/template"
	"strconv"
)

func ParseJson(fileContents []byte, tmpl *template.TemplateWithDetails) error {
	elements, err := parse(fileContents)
	if err != nil {
		return err
	}
	*tmpl = template.TemplateWithDetails{
		AWSTemplateFormatVersion: elements["AWSTemplateFormatVersion"],
		Description:              elements["Description"],
		Metadata:                 elements["Metadata"],
		Parameters:               elements["Parameters"],
		Mappings:                 elements["Mappings"],
		Conditions:               elements["Conditions"],
		Transform:                elements["Transform"],
		Resources:                elements["Resources"],
		Outputs:                  elements["Outputs"]}
	return nil
}

func parse(fileContents []byte) (map[string]*template.TemplateElement, error) {
	elements := make(map[string]*template.TemplateElement)
	err := jsonparser.ObjectEach(fileContents, nestedJsonObjectIterator(string(fileContents), 0, elements))
	if err != nil {
		return nil, err
	}
	return elements, nil
}

func nestedJsonObjectIterator(fileContents string, offset int, parentChildrenMap map[string]*template.TemplateElement) func(key []byte, value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, valueStartOffset int) error {
	var iterator func(key []byte, value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, valueStartOffset int) error

	iterator = func(key []byte, value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, valueStartOffset int) error {
		line, character := lineAndCharacter(fileContents, startOffset+offset)
		element := &template.TemplateElement{
			Name:   string(key),
			Line:   line,
			Column: character,
			Type:   translateElementType(dataType)}
		parentChildrenMap[element.Name] = element
		if dataType == jsonparser.Object {
			elementChildrenMap := make(map[string]*template.TemplateElement)
			element.Children = elementChildrenMap
			return jsonparser.ObjectEach(value, nestedJsonObjectIterator(fileContents, offset+valueStartOffset, elementChildrenMap))
		} else if dataType == jsonparser.Array {
			elementChildrenArray := make([]*template.TemplateElement, 0)
			_, err := jsonparser.ArrayEach(value, nestedJsonArrayIterator(fileContents, offset+valueStartOffset, &elementChildrenArray))
			element.Children = &elementChildrenArray
			return err
		} else {
			element.Value = parsePrimitiveValue(value, dataType)
		}
		return nil
	}

	return iterator
}

func nestedJsonArrayIterator(fileContents string, offset int, parentChildrenArray *[]*template.TemplateElement) func(value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, err error) {
	var iterator func(value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, err error)
	iterator = func(value []byte, dataType jsonparser.ValueType, startOffset int, endOffset int, err error) {
		line, character := lineAndCharacter(string(fileContents), startOffset+offset)
		element := &template.TemplateElement{
			Name:   "[" + strconv.Itoa(len(*parentChildrenArray)) + "]",
			Line:   line,
			Column: character,
			Type:   translateElementType(dataType)}
		*parentChildrenArray = append(*parentChildrenArray, element)
		if dataType == jsonparser.Object {
			elementChildrenMap := make(map[string]*template.TemplateElement)
			element.Children = elementChildrenMap
			jsonparser.ObjectEach(value, nestedJsonObjectIterator(fileContents, offset+startOffset, elementChildrenMap))
		} else if dataType == jsonparser.Array {
			elementChildrenArray := make([]*template.TemplateElement, 1)
			jsonparser.ArrayEach(value, nestedJsonArrayIterator(fileContents, offset+startOffset, &elementChildrenArray))
			element.Children = &elementChildrenArray
		} else {
			element.Value = parsePrimitiveValue(value, dataType)
		}
	}
	return iterator
}

func lineAndCharacter(input string, offset int) (line int, character int) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0
	}

	line = 1

	for i, b := range input {
		if b == lf {
			if i < offset {
				line++
				character = 0
			}
		} else {
			character++
		}
		if i == offset {
			break
		}
	}

	return line, character
}

func translateElementType(valueType jsonparser.ValueType) template.TemplateElementValueType {
	switch valueType {
	case jsonparser.NotExist:
		return template.NotExist
	case jsonparser.String:
		return template.String
	case jsonparser.Number:
		return template.Number
	case jsonparser.Object:
		return template.Object
	case jsonparser.Array:
		return template.Array
	case jsonparser.Boolean:
		return template.Boolean
	case jsonparser.Null:
		return template.Null
	default:
		return template.Unknown
	}
}

func parsePrimitiveValue(value []byte, valueType jsonparser.ValueType) interface{} {
	switch valueType {
	case jsonparser.String:
		return string(value)
	case jsonparser.Number:
		var val interface{}
		val, err := jsonparser.ParseInt(value)
		if err != nil {
			val, _ = jsonparser.ParseFloat(value)
		}
		return val
	case jsonparser.Boolean:
		val, _ := jsonparser.ParseBoolean(value)
		return val
	case jsonparser.Null:
		return nil
	default:
		return nil
	}
}
