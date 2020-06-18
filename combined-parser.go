package logparser

import (
	"fmt"
	"time"
)

type CombinedParser struct {
	Parsers []LogFileParser
}

func (c *CombinedParser) Process(start, end time.Time) ([]LogLine, error) {
	logLines := make([]LogLine, 0)

	for _, parser := range c.Parsers {
		localLogLines, err := parser.Process(start, end)
		if err != nil {
			fmt.Println("ERROR:", err)
		}

		logLines = append(logLines, localLogLines...)
	}

	return logLines, nil
}
