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

package helpers

// SliceContains checks if slice contains given string.
func SliceContains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsPlainMap checks if map is plain. Plain map means that it's non-nil and it doesn't contain nested maps.
//It's used in checkWhereIsNil().
func IsPlainMap(mp map[string]interface{}) bool {
	// First we check is it more complex. If so - it is worth investigating and we should stop checking.
	for _, m := range mp {
		if _, ok := m.(map[string]interface{}); ok {
			return false
		} else if _, ok := m.([]interface{}); ok {
			return false
		}
	}
	// Ok, it isn't. So is there any <nil>?
	if mapContainsNil(mp) { // Yes, it is - so it is a map worth investigating. This is not the map we're looking for.
		return false
	}

	return true // There is no <nil> and no complexity - it is a plain, non-nil map.
}

// IsPlainSlice checks if slice is plain. Slice is plain if it's non-nil and doesn't contain nested maps.
//It's used in checkWhereIsNil().
func IsPlainSlice(slc []interface{}) bool {
	// The same flow as in `isPlainMap` function.
	for _, s := range slc {
		if _, ok := s.(map[string]interface{}); ok {
			return false
		} else if _, ok := s.([]interface{}); ok {
			return false
		}
	}

	if sliceContainsNil(slc) {
		return false
	}

	return true
}

// Discard looks for elements which are not the same and return only unique.
func Discard(slice []interface{}, n interface{}) []interface{} {
	result := []interface{}{}
	for _, s := range slice {
		if s != n {
			result = append(result, s)
		}
	}
	return result
}

// IsNonStringFloatBool checks if the element is non-string, non-float64, non-boolean. Then it is another node or <nil>. There is no other option.
func IsNonStringFloatBool(v interface{}) bool {
	var isString, isFloat, isBool bool
	if _, ok := v.(string); ok {
		isString = true
	} else if _, ok := v.(float64); ok {
		isFloat = true
	} else if _, ok := v.(bool); ok {
		isBool = true
	}
	if !isString && !isFloat && !isBool {
		return true
	}
	return false
}

func mapContainsNil(mp map[string]interface{}) bool {
	for _, m := range mp {
		if m == nil {
			return true
		}
	}
	return false
}

func sliceContainsNil(slice []interface{}) bool {
	for _, s := range slice {
		if s == nil {
			return true
		}
	}
	return false
}
