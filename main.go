package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat"
)

type MetricSummary struct {
	Mean   float64
	Median float64
	Min    float64
	Max    float64
	Std    float64
}

type FileSummary struct {
	StartTime string
	EndTime   string
	Sport     string
	Filename  string
	Duration  float64
	Distance  float64
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
			FloatToString(summary.Distance),
			FloatToString(summary.Heartrate.Median),
			FloatToString(summary.Heartrate.Min),
			FloatToString(summary.Heartrate.Max),
			FloatToString(summary.Heartrate.Std),
			FloatToString(summary.Elevation.Min),
			FloatToString(summary.Elevation.Max),
			FloatToString(summary.Cadence.Median),
			FloatToString(summary.Cadence.Min),
			FloatToString(summary.Cadence.Max),
			FloatToString(summary.Cadence.Std),
			FloatToString(summary.Power.Median),
			FloatToString(summary.Power.Min),
			FloatToString(summary.Power.Max),
			FloatToString(summary.Power.Std),
		}

		writer.Write(entry)
	}
}

// haversine function calculates the distance between two points on the Earth
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters

	if lat1 == 0.0 || lon1 == 0.0 || lat2 == 0.0 || lon2 == 0.0 {
		return 0.0
	}

	var φ1 = lat1 * math.Pi / 180
	var φ2 = lat2 * math.Pi / 180
	var Δφ = (lat2 - lat1) * math.Pi / 180
	var Δλ = (lon2 - lon1) * math.Pi / 180

	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // in meters
}

func getMetricSummary(data []float64) MetricSummary {
	sort.Float64s(data)
	median := stat.Quantile(0.5, stat.Empirical, data, nil)

	return MetricSummary{
		Mean:   stat.Mean(data, nil),
		Median: median,
		Min:    slices.Min(data),
		Max:    slices.Max(data),
		Std:    stat.StdDev(data, nil),
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
	var prevLat, prevLong float64
	var heartrates, eleveations, cadences, powers []float64

	for _, record := range records {
		curTime := record[0]
		sport := record[1]
		filename := record[2]

		lat, _ := strconv.ParseFloat(record[3], 64)
		long, _ := strconv.ParseFloat(record[4], 64)

		elevation, _ := strconv.ParseFloat(record[5], 64)
		cadence, _ := strconv.ParseFloat(record[6], 64)
		hr, _ := strconv.ParseFloat(record[7], 64)
		power, _ := strconv.ParseFloat(record[8], 64)

		if cFilename == "" {
			cFilename = filename
		} else if cFilename != filename {
			summary := summaries[cFilename]

			summary.Heartrate = getMetricSummary(heartrates)
			summary.Elevation = getMetricSummary(eleveations)
			summary.Cadence = getMetricSummary(cadences)
			summary.Power = getMetricSummary(powers)

			summaries[cFilename] = summary

			cFilename = filename
		}

		if summary, ok := summaries[cFilename]; ok {
			startTime, _ := time.Parse(time.RFC3339, summary.StartTime)
			endTime, _ := time.Parse(time.RFC3339, curTime)
			distance := haversine(prevLat, prevLong, lat, long)

			summary.EndTime = curTime
			summary.Duration = endTime.Sub(startTime).Seconds()
			summary.Distance = summary.Distance + distance

			summaries[cFilename] = summary
		} else {
			heartrates, eleveations, cadences, powers =
				[]float64{}, []float64{}, []float64{}, []float64{}

			summaries[cFilename] = FileSummary{
				StartTime: curTime,
				Sport:     sport,
				Filename:  cFilename,
				Distance:  0.0,
			}
		}

		prevLat = lat
		prevLong = long

		heartrates = append(heartrates, hr)
		eleveations = append(eleveations, elevation)
		cadences = append(cadences, cadence)
		powers = append(powers, power)
	}

	WriteCSV(summaries, outputFile)
}
