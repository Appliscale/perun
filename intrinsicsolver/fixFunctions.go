package intrinsicsolver

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// FixFunctions : takes []byte file and firstly converts all single quotation marks to double ones (anything between single ones is treated as the rune in GoLang), then deconstructs file into lines, checks for intrinsic functions and converts them to equivalent JSON representation and saves the result to temporary file, then opens it and returns []byte array
func FixFunctions(template []byte) []byte {
	var quotationProcessed, temporaryResult []string
	//parse original file into lines
	preLines := parseFileIntoLines(template)
	//fix single-double quotation marks issue
	for _, e := range preLines {
		var fixed string
		if strings.Contains(e, "'") {
			fixed = strings.Replace(e, "'", "\"", -1)
		} else {
			fixed = e
		}

		quotationProcessed = append(quotationProcessed, fixed)
	}

	lines := quotationProcessed

	//check line by line
	for idx, d := range lines {
		//store the actual line in cache
		cache := d
		var postCache string

		startCount := 1 //set the inital guess for intrinsic functions in line to one
		var topCount int
		if strings.Contains(d, "!") {
			topCount = strings.Count(d, "!") //count the actual number of that functions
		} else if strings.Contains(d, "Fn::Base64") { //check for userData extended Fn::Base64 tag
			topCount = strings.Count(d, "Fn::Base64")
		}

		for startCount <= topCount { //while there are any intrinsic functions inside the line

			if startCount > 1 {
				cache = postCache //update the after-state cache
			}

			//first fix the most complex tags
			fixSub(&cache, lines, idx)
			fixIf(&cache, lines, idx)

			//then all the rest
			fixSelect(&cache, lines, idx)
			fixGetAZs(&cache)
			fixEquals(&cache)
			fixGetAtt(&cache)
			fixImportValue(&cache)
			fixFindInMap(&cache)
			fixRef(&cache)
			fixUserData(&cache, lines, idx)
			fixSplit(&cache)

			//iteration-over line becomes the line from that iteration
			postCache = cache
			//if any of the functions were resolved, we are done with that one
			startCount = startCount + 1

		}
		//append cached line to the slice of temporary results - if there were no fixes, cache is 'd' as in the 'check by line' iteration entry, if there were - cache is the fixed line
		temporaryResult = append(temporaryResult, cache)

	}

	var counter int
	//function DeleteFromSlice originaly were leaving the empty lines in the end for every deleted line in result so it was exchanged with leaving the "%%%" mark instead to easily cut the unnecessary lines in the end
	for i := len(temporaryResult) - 1; i >= 0; i-- {
		if strings.Contains(temporaryResult[i], "%%%") {
			counter++ //counter counts how many are lines to delete
		}
	}

	newLen := len(temporaryResult) - counter //new length is the length of a temporary file minus how many lines to be deleted

	//writeLines saves the processed result to file (if there would be any errors, it could be investigated there)
	if err := writeLines(temporaryResult, "preprocessed.yml", newLen); err != nil {
		log.Fatalf("writeLines: %s", err)
	}

	//then the temporary result is opened...
	preprocessedTemplate, err := ioutil.ReadFile("preprocessed.yml")
	if err != nil {
		fmt.Println(err)
	}

	//... and returned as []byte
	return preprocessedTemplate
}

//ADDITIONAL FUNCTIONS

//from:
//https://stackoverflow.com/questions/39862613/how-to-split-multiple-delimiter-in-golang
//splitter splits the string by multiple delimiters
func splitter(s string, splits string) []string {
	m := make(map[rune]int)
	for _, r := range splits {
		m[r] = 1
	}

	runeSplitter := func(r rune) bool {
		return m[r] == 1
	}

	return strings.FieldsFunc(s, runeSplitter)
}

//parseFileIntoLines is reading the []byte file and returns it line by line as []string slice
func parseFileIntoLines(template []byte) []string {
	bytesReader := bytes.NewReader(template)
	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytesReader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

//from:
//https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
//writeLines takes []string slice and writes it element by element as line by line file
func writeLines(lines []string, path string, c int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines[:c] {
		fmt.Fprintln(w, line)
	}
	fmt.Fprintln(w, "")

	return w.Flush()
}

//deletes element from slice and places deletion mark (%%%) at the end (marks are then used by function FixFunctions to compute the final file length) - probably to be fixed by better function
func deleteFromSlice(slice []string, index int, n int) {
	copy(slice[index+n:], slice[index+1+n:])
	slice[len(slice)-1] = "%%%"
	slice = slice[:len(slice)-1]
}

//from:
//https://rosettacode.org/wiki/Determine_if_a_string_is_numeric#Go
//check if string is numeric
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
