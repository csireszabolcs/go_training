package logparser

import (
	"fmt"
	"time"
)

type LogLine struct {
	TimeStamp time.Time
	Text      string
	Filename  string
}

func (l LogLine) String() string {
	return fmt.Sprintf("(%v) : [%v] %v", l.Filename, l.TimeStamp, l.Text)
}

type LogFileParser interface {
	Process(start, end time.Time) ([]LogLine, error)
}
