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

// Package intrinsicsolver contains of its own solver.
package intrinsicsolver

// Function indentations checks how much an element is indented by counting all the spaces encountered in searching for the first non-space character in line.
func indentations(line string) int {
	var i int
	for string(line[i]) == " " {
		i++
	}
	return i
}
