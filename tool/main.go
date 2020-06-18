package main

import (
	"flag"
	"fmt"
	"logparser"
	"os"
	"time"
)

type jsonLogLine struct {
	Ts string `json:"ts"`
}

func main() {
	//startStringRef := flag.String("start", "2020-06-11T13:48:00.000", "The start time to look for in logs")
	//endStringRef := flag.String("end", "2020-06-11T13:59:00.000", "The end time to look for in logs")
	startStringRef := flag.String("start", "2020-06-12T00:26:00.000", "The start time to look for in logs")
	endStringRef := flag.String("end", "2020-06-12T00:50:00.000", "The end time to look for in logs")
	durationRef := flag.String("dur", "", "The duration from start, e.g: [1h, 1m, 1s]")

	flag.Parse()

	// Check here because I'm paranoid and these are user supplied pointers that will be dereferenced
	if startStringRef == nil || (endStringRef == nil && durationRef == nil) {
		fmt.Println("Missing start and / or end time or duration")
		os.Exit(-1)
	}

	parser := logparser.CombinedParser{
		Parsers: []logparser.LogFileParser{
			//&logparser.HttpAccessLineParser{"C:/tmp/logs/logs/httpd/access.*.log"},
			//&logparser.JsonLogLineParser{"C:/tmp/logs/logs/vizqlserver/nativeapi_vizqlserver*.txt"},
			//&logparser.VizportalLogLineParser{"C:/tmp/logs/logs/vizportal/vizportal_node*"},
			&logparser.RedisLogLineParser{"C:/tmp/logs/logs/cacheserver/redis_*.log"},
		},
	}

	var start time.Time
	var end time.Time
	var err error

	// Check times validity
	//start, err := time.Parse("2006-01-02T15:04:05", "2020-06-11T13:48:00.000")
	start, err = time.Parse("2006-01-02T15:04:05", *startStringRef)
	if err != nil {
		fmt.Println("Malformed start time:", err)
		os.Exit(-1)
	}


	if len(*durationRef) > 0 {
		fmt.Println("Using duration arg", *durationRef)
		duration, err:= time.ParseDuration(*durationRef)
		if err != nil {
			fmt.Println("Malformed duration:", err)
			os.Exit(-1)
		}
		end = start.Add(duration)
		fmt.Println("End date from duration", end)
	} else {
		fmt.Println("Using end date", *endStringRef)
		end, err = time.Parse("2006-01-02T15:04:05", *endStringRef)
		if err != nil {
			fmt.Println("Malformed end time:", err)
			os.Exit(-1)
		}
	}


	logLines, err := parser.Process(start, end)

	if err != nil {
		panic(err)
	}

	for _, line := range logLines {
		fmt.Println(line)
	}

}
