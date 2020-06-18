package logparser

import (
	"encoding/json"
	"fmt"
	"time"
)

type JsonLogLineParser struct {
	FileGlob string
}

func (h *JsonLogLineParser) Process(start, end time.Time) ([]LogLine, error) {

	// storage for output
	logLinesMatched := make([]LogLine, 0)

	//
	err := forEachMatchedFile(h.FileGlob, func(match string) error {
		// find all lines
		logLines, err := h.parseFile(match)
		if err != nil {
			return fmt.Errorf("while trying to parse log lines from file '%v': %v", match, err)
		}

		// filter the lines by timestamp
		matchedLines := filterLogLines(start, end, logLines)

		// append the filtered lines to the matched lines
		logLinesMatched = append(logLinesMatched, matchedLines...)

		return nil
	})

	// return the error as-is
	return logLinesMatched, err

}

type jsonLogLine struct {
	Ts string `json:"ts"`
}

func (h *JsonLogLineParser) parseFile(filename string) ([]LogLine, error) {
	// create the empty container
	logLines := make([]LogLine, 0)

	err := forEachLineOfFile(filename, func(lineText string) error {

		var jsonLine jsonLogLine

		// parse json
		if err := json.Unmarshal([]byte(lineText), &jsonLine); err != nil {
			return fmt.Errorf("while attempting to parse json line '%v': %v", lineText, err)
		}

		// parse & log errors (for now)
		parsedTime, err := time.Parse("2006-01-02T15:04:05", jsonLine.Ts)
		if err != nil {
			return fmt.Errorf("while parsing time'%v': %v", jsonLine.Ts, err)
		}

		logLines = append(logLines, LogLine{
			TimeStamp: parsedTime,
			Text:      lineText,
			Filename:  filename,
		})

		return nil
	})

	return logLines, err
}
