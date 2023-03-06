package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type inputType string

var (
	logInput        inputType = "log"
	filestreamInput           = "fs"
)

type result struct {
	Run        inputType `json:"run"`
	Start      time3339  `json:"start"`
	End        time3339  `json:"end"`
	Duration   float64   `json:"duration_sec"`
	Throughput float64   `json:"throughput_MB_sec"`
	Bytes      float64   `json:"bytes"`
	Lines      float64   `json:"lines"`
}

type report struct {
	Result1         result
	Result2         result
	DurationDelta   float64 `json:"duration_delta_sec"`
	BytesDelta      float64 `json:"bytes_delta"`
	ThroughputDelta float64 `json:"throughput_delta_MB_sec"`
}

type time3339 struct {
	time.Time
}

func (m *time3339) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}
	var format string

	if strings.LastIndexByte(str, 'Z') != -1 {
		format = `"` + time.RFC3339Nano + `"`
	} else {
		format = `"2006-01-02T15:04:05.999-0700"`
	}

	tt, err := time.Parse(format, string(data))
	*m = time3339{tt}
	return err
}

func (m *time3339) ToTime() time.Time {
	return m.Time
}

func main() {
	if len(os.Args) != 3 { // first is the program name
		fmt.Println(`
Result files are not set.

Usage:
go run main.go result1.json result2.json`)
		os.Exit(1)
		return
	}

	result1Filename := os.Args[1]
	result2Filename := os.Args[2]

	result1, err := getResult(result1Filename)
	if err != nil {
		log.Fatal(fmt.Errorf("failed get results from %s file: %w", result1Filename, err))
	}
	result2, err := getResult(result2Filename)
	if err != nil {
		log.Fatal(fmt.Errorf("failed get results from %s file: %w", result2Filename, err))
	}

	report := compare(result1, result2)
	reportBytes, err := json.Marshal(report)
	if err != nil {
		log.Fatal(fmt.Errorf("failed get serialise the report: %w", err))
	}
	fmt.Println(string(reportBytes))
}

func delta(x, y float64) float64 {
	return (1 - math.Min(x, y)/math.Max(x, y)) * 100
}

func compare(r1, r2 result) report {
	return report{
		Result1:         r1,
		Result2:         r2,
		DurationDelta:   delta(r1.Duration, r2.Duration),
		BytesDelta:      delta(r1.Bytes, r2.Bytes),
		ThroughputDelta: delta(r1.Throughput, r2.Throughput),
	}
}

func getResult(filename string) (r result, err error) {
	resultBytes, err := os.ReadFile(filename)
	if err != nil {
		return r, fmt.Errorf("failed to read the result file: %w", err)
	}
	err = json.Unmarshal(resultBytes, &r)
	if err != nil {
		return r, fmt.Errorf("failed to parse the result file: %w", err)
	}

	r.Duration = r.End.Sub(r.Start.ToTime()).Seconds()
	r.Throughput = (r.Bytes / r.Duration) / 1024 / 1024 // MB/sec

	return r, nil
}
