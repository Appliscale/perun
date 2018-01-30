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

// Package utilities provides various helpers used
// here and there over the code base.
package utilities

import (
	"fmt"
	"time"
)

const Motto = "Swiss army knife for AWS CloudFormation templates"

const ReleaseName = "Nimbostratus"
const VersionNumber = "1.1.0-beta"

func VersionStatus() string {
	return fmt.Sprintf("       perun %s (%s release) - %s", VersionNumber, ReleaseName, Motto)
}

func TruncateDuration(d time.Duration) time.Duration {
	return -(d - d%(time.Duration(1)*time.Second))
}
