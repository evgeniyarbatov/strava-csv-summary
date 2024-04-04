package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type MetricSummary struct {
	Mean   float64
	Median float64
	Min    float64
	Max    float64
	Std    float64
	Count  int
}

type FileSummary struct {
	StartTime string
	EndTime   string
	Sport     string
	Filename  string
	Duration  float64
	Distance  int
	Heartrate MetricSummary
	Elevation MetricSummary
	Cadence   MetricSummary
	Power     MetricSummary
}

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func WriteCSV(
	summaries map[string]FileSummary,
	filePath string,
) {
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for _, summary := range summaries {
		entry := []string{
			summary.StartTime,
			summary.EndTime,
			summary.Sport,
			summary.Filename,
			FloatToString(summary.Duration),
			strconv.Itoa(summary.Distance),
		}

		writer.Write(entry)
	}
}

func main() {
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	summaries := make(map[string]FileSummary)

	var cFilename string

	for _, record := range records {
		curTime := record[0]
		sport := record[1]
		filename := record[2]

		if cFilename == "" || cFilename != filename {
			cFilename = filename
		}

		if summary, ok := summaries[cFilename]; ok {
			startTime, _ := time.Parse(time.RFC3339, summary.StartTime)
			endTime, _ := time.Parse(time.RFC3339, curTime)

			summary.Duration = endTime.Sub(startTime).Seconds()
			summary.EndTime = curTime

			summaries[cFilename] = summary
		} else {
			summaries[cFilename] = FileSummary{
				StartTime: curTime,
				Sport:     sport,
				Filename:  cFilename,
			}
		}
	}

	WriteCSV(summaries, outputFile)
}
