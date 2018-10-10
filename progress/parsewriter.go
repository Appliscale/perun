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

package progress

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const createCompleteStatus = "CREATE_COMPLETE"
const createInProgressStatus = "CREATE_IN_PROGRESS"
const createFailedStatus = "CREATE_FAILED"
const deleteCompleteStatus = "DELETE_COMPLETE"
const deleteFailedStatus = "DELETE_FAILED"
const deleteInProgressStatus = "DELETE_IN_PROGRESS"
const reviewInProgressStatus = "REVIEW_IN_PROGRESS"
const rollbackCompleteStatus = "ROLLBACK_COMPLETE"
const rollbackFailedStatus = "ROLLBACK_FAILED"
const rollbackInProgressStatus = "ROLLBACK_IN_PROGRESS"
const updateCompleteStatus = "UPDATE_COMPLETE"
const updateCompleteCleanupInProgressStatus = "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"
const updateInProgressStatus = "UPDATE_IN_PROGRESS"
const updateRollbackCompleteStatus = "UPDATE_ROLLBACK_COMPLETE"
const updateRollbackCompleteCleanupInProgressStatus = "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"
const updateRollbackFailedStatus = "UPDATE_ROLLBACK_FAILED"
const updateRollbackInProgressStatus = "UPDATE_ROLLBACK_IN_PROGRESS"

const add = "Add"
const remove = "Remove"
const modify = "Modify"

// ParseWriter structure describes functions to stain text.
type ParseWriter struct {
	linesPrinted   int
	bgRed          func(a ...interface{}) string
	fgRed          func(a ...interface{}) string
	fgOrange       func(a ...interface{}) string
	bgOrange       func(a ...interface{}) string
	grey           func(a ...interface{}) string
	bgGreen        func(a ...interface{}) string
	fgGreen        func(a ...interface{}) string
	cyan           func(a ...interface{}) string
	statusColorMap map[string]func(a ...interface{}) string
}

// NewParseWriter is used to show progress in colorful table.
func NewParseWriter() (pw *ParseWriter) {
	pw = &ParseWriter{}
	pw.linesPrinted = 0
	pw.bgRed = color.New(color.BgHiRed).SprintFunc()
	pw.fgRed = color.New(color.FgRed).SprintFunc()
	pw.fgOrange = color.New(color.FgHiYellow).SprintFunc()
	pw.bgOrange = color.New(color.BgHiYellow).SprintFunc()
	pw.grey = color.New(color.FgHiWhite).SprintFunc()
	pw.bgGreen = color.New(color.BgGreen).SprintFunc()
	pw.fgGreen = color.New(color.FgHiGreen).SprintFunc()
	pw.cyan = color.New(color.FgCyan).SprintFunc()

	pw.statusColorMap = map[string]func(a ...interface{}) string{
		createFailedStatus:                            pw.bgRed,
		rollbackFailedStatus:                          pw.bgRed,
		rollbackCompleteStatus:                        pw.fgRed,
		updateRollbackCompleteStatus:                  pw.fgRed,
		updateRollbackInProgressStatus:                pw.fgRed,
		rollbackInProgressStatus:                      pw.fgRed,
		deleteFailedStatus:                            pw.bgRed,
		updateRollbackFailedStatus:                    pw.bgRed,
		deleteCompleteStatus:                          pw.grey,
		createInProgressStatus:                        pw.fgOrange,
		updateRollbackCompleteCleanupInProgressStatus: pw.bgOrange,
		deleteInProgressStatus:                        pw.fgOrange,
		updateCompleteCleanupInProgressStatus:         pw.fgOrange,
		updateInProgressStatus:                        pw.fgOrange,
		createCompleteStatus:                          pw.bgGreen,
		updateCompleteStatus:                          pw.bgGreen,
		reviewInProgressStatus:                        pw.cyan,
		add:    pw.bgGreen,
		remove: pw.bgRed,
		modify: pw.fgOrange,
	}
	return
}

// Write prints lines.
func (pw *ParseWriter) Write(p []byte) (n int, err error) {
	var newString = pw.colorStatuses(string(p))
	fmt.Print(newString)
	pw.linesPrinted += strings.Count(newString, "\n") - strings.Count(newString, "\033[A")
	return len(p), nil
}
func (pw *ParseWriter) colorStatuses(s string) string {
	for status, colorizeFun := range pw.statusColorMap {
		if strings.Contains(s, status) {
			s = strings.Replace(s, status, colorizeFun(status), -1)
		}
	}
	return s
}

func (pw *ParseWriter) returnWritten() {
	for i := 0; i < pw.linesPrinted; i++ {
		fmt.Print("\033[A")
	}
	pw.linesPrinted = 0
	return
}
