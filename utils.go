package logparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Runs `fn(filepath)` for each file matching the Glob `glob`.
// Errors are gathered in an slice, because they are to be reported, but not causes of failiure
func forEachMatchedFile(glob string, fn func(string) error) error {
	// find relevant files
	matches, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("while trying to glob '%v': %v", glob, err)
	}

	fmt.Println("Matches", matches)
	var errors []error

	// process each file using the user-supplied fn
	for _, match := range matches {
		fmt.Println("file: ", match)

		// add the error to the errors list if needed
		if err := fn(match); err != nil {
			errors = append(errors, fmt.Errorf("while processing file '%v': %v", err))
		}
	}

	logErrors(errors)

	return nil
}

func logErrors(errors []error) {
	// log each error
	for _, err := range errors {
		fmt.Println("ERROR:", err)
	}
}

// Helper that filters the list of log lines to include only entries between start and end-time
func filterLogLines(start, end time.Time, logLines []LogLine) []LogLine {

	// storage for output
	logLinesMatched := make([]LogLine, 0)

	// do the checking
	for _, line := range logLines {
		if checkIfTimeInBetween(start, end, line.TimeStamp) {
			logLinesMatched = append(logLinesMatched, line)
		}
	}

	return logLinesMatched
}

// Execute `fn(line)` for each of the lines in the file
func forEachLineOfFile(filename string, fn func(line string) error) error {

	// open the file
	file, err := os.Open(filename) // For read access.
	if err != nil {
		return fmt.Errorf("while opening '%v' for reading: %v", filename, err)
	}
	// close the file
	defer file.Close()

	var errors []error

	// read line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// read line
		lineText := scanner.Text()
		// call fn
		if err := fn(lineText); err != nil {
			// append error
			errors = append(errors, err)
		}
	}

	logErrors(errors)

	return nil
}

// Checks if eventTime is between start and end
func checkIfTimeInBetween(start, end, eventTime time.Time) bool {
	return eventTime.After(start) && eventTime.Before(end)
}