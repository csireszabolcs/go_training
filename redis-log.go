package logparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type RedisLogLineParser struct {
	FileGlob string
}

var RedisLineRe = regexp.MustCompile(`^.* (\d+ \w+ \d+:\d+:\d+)`)

func ParseRedisTime(redisTime [][]byte) (time.Time, error) {
	// parse the time
	currentTime :=  time.Now()
	currentYear := strconv.Itoa(currentTime.Year())

	// Is there a cleaner solution for time without year
	// e.g: [11176] 10 Jun 17:20:11.311 * DB saved on disk
	parsedTime, err := time.Parse("2006 02 Jan 15:04:05", currentYear + " " + string(redisTime[1]))
	if err != nil {
		return time.Time{}, fmt.Errorf("while parsing time: %v", err)
	}

	return parsedTime, nil
}

func (h *RedisLogLineParser) parseFile(filename string) ([]LogLine, error) {

	// create the empty container
	logLines := make([]LogLine, 0)

	err := forEachLineOfFile(filename, func(lineText string) error {
		// skipping logo lines like: //
		//          _.-``__ ''-._
		//     _.-``    `.  `_.  ''-._           Redis 3.0.503 (00000000/0) 64 bit
		// .-`` .-```.  ```\/    _.,_ ''-._
		//(    '      ,       .-`  | `,    )     Running in standalone mode
		//|`-._`-...-` __...-.``-._|'` _.-'|     Port: 8982
		//|    `-._   `._    /     _.-'    |     PID: 5080
		// `-._    `-._  `-./  _.-'    _.-'
		//|`-._`-._    `-.__.-'    _.-'_.-'|
		//|    `-._`-._        _.-'_.-'    |           http://redis.io
		// `-._    `-._`-.__.-'_.-'    _.-'
		//|`-._`-._    `-.__.-'    _.-'_.-'|
		//|    `-._`-._        _.-'_.-'    |
		// `-._    `-._`-.__.-'_.-'    _.-'
		//     `-._    `-.__.-'    _.-'
		//         `-._        _.-'
		//             `-.__.-'

		result := RedisLineRe.FindSubmatch([]byte(lineText))
		if result != nil {
			lineParsedTime, err := ParseRedisTime(result)
			if err != nil {
				return err
			}

			logLines = append(logLines, LogLine{
				TimeStamp: lineParsedTime,
				Text:      lineText,
				Filename:  filename,
			})
		}

		return nil
	})

	return logLines, err
}

func (h *RedisLogLineParser) Process(start, end time.Time) ([]LogLine, error) {

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
