package logparser

import (
	"fmt"
	"regexp"
	"time"
)

var httpLogLineRe = regexp.MustCompile(`^.* - (\d+.\d+.\d+T\d+:\d+:\d+\.\d+)`)

func ParseHttpLine(line string) (time.Time, error) {

	result := httpLogLineRe.FindSubmatch([]byte(line))

	// if we have errors
	if result == nil {
		return time.Time{}, fmt.Errorf("cannot find the timestamp in the line '%v'", line)
	}

	// parse the time
	parsedTime, err := time.Parse("2006-01-02T15:04:05", string(result[1]))
	if err != nil {
		return time.Time{}, fmt.Errorf("while parsing time: %v", err)
	}

	return parsedTime, nil
}

type HttpAccessLineParser struct {
	FileGlob string
}

func (h *HttpAccessLineParser) parseFile(filename string) ([]LogLine, error) {

	// create the empty container
	logLines := make([]LogLine, 0)

	err := forEachLineOfFile(filename, func(lineText string) error {
		lineParsedTime, err := ParseHttpLine(lineText)
		if err != nil {
			return err
		}

		logLines = append(logLines, LogLine{
			TimeStamp: lineParsedTime,
			Text:      lineText,
			Filename:  filename,
		})

		return nil
	})

	return logLines, err
}

func (h *HttpAccessLineParser) Process(start, end time.Time) ([]LogLine, error) {

	// storage for output
	logLinesMatched := make([]LogLine, 0)

	// check each file
	err := forEachMatchedFile(h.FileGlob, func(match string) error {

		// find all lines
		logLines, err := h.parseFile(match)
		//logLines, err := CheckHttpFileForLines(match)
		if err != nil {
			return fmt.Errorf("while trying to parse log lines from file '%v': %v", match, err)
		}

		// filter the lines by timestamp
		matchedLines := filterLogLines(start, end, logLines)

		// append the filtered lines to the matched lines
		logLinesMatched = append(logLinesMatched, matchedLines...)

		return nil
	})

	return logLinesMatched, err
}
