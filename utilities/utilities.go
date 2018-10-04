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

// Package utilities provides various helpers used
// here and there over the code base.
package utilities

import (
	"fmt"
	"os"
	"time"
)

// Motto - perun's motto.
const Motto = "Swiss army knife for AWS CloudFormation templates"

// ReleaseName - name of the release.
const ReleaseName = "Altostratus"

// VersionNumber - number of the release.
const VersionNumber = "1.3.1"

// VersionStatus shows perun's release.
func VersionStatus() string {
	return fmt.Sprintf("perun %s (%s release) - %s", VersionNumber, ReleaseName, Motto)
}

// TruncateDuration prepares shorter message with duration.
func TruncateDuration(d time.Duration) time.Duration {
	return -(d - d%(time.Duration(1)*time.Second))
}

// CheckErrorCodeAndExit checks if error exists.
func CheckErrorCodeAndExit(err error) {
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// CheckFlagAndExit checks error flag.
func CheckFlagAndExit(valid bool) {
	if valid {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
