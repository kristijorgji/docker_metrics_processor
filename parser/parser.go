package parser

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"../models"
)

const mysqlDateFormat = "2006-01-02 15:04:05.000000"

// Parse the file to the proper data structure
func Parse(filePath string) []models.ServiceMetrics {
	readFile, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var dateTimePattern = regexp.MustCompile("(\\d{4}-\\d{2}-\\d{2}.{15})")
	var dataPattern = regexp.MustCompile("(?P<containerId>.{12}) {2,}(?P<containerName>[^ ]+) {2,}(?P<cpuPercentage>[^ ]+)% {2,}(?P<memoryUsage>[^ ]+) / (?P<memoryLimit>[^ ]+) {2,}(?P<memoryPercentage>[^ ]+)% {2,}")

	var metrics []models.ServiceMetrics
	var serviceMetrics *models.ServiceMetrics = new(models.ServiceMetrics)

	var i = 0
	var datetime string
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if dateTimePattern.MatchString(line) {
			date := dateTimePattern.FindAllString(line, -1)[0]
			t, err := time.Parse(time.RFC3339, date)
			if err != nil {
				log.Fatalf("failed to parse date: %s", err)
			}
			datetime = t.Format(mysqlDateFormat)
		} else {
			match := dataPattern.FindStringSubmatch(line)
			if len(match) == 0 {
				continue
			}

			result := make(map[string]string)
			for i, name := range dataPattern.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			if result["containerId"] == "CONTAINER ID" {
				continue
			}

			serviceMetrics.Datetime = datetime
			serviceMetrics.ContainerID = result["containerId"]
			serviceMetrics.ContainerName = result["containerName"]

			tmp, _ := strconv.ParseFloat(result["cpuPercentage"], 32)
			serviceMetrics.CPUPercentage = float32(tmp)

			serviceMetrics.MemoryUsageInMib = toMib(result["memoryUsage"])
			serviceMetrics.MemoryLimitInMib = toMib(result["memoryLimit"])

			tmp, _ = strconv.ParseFloat(result["memoryPercentage"], 32)
			serviceMetrics.MemoryPercentage = float32(tmp)

			metrics = append(metrics, *serviceMetrics)
			serviceMetrics = new(models.ServiceMetrics)
		}

		i++
	}

	readFile.Close()

	return metrics
}

var sizePattern = regexp.MustCompile("([\\d\\.]+)(\\w+)$")

const (
	b   = "B"
	mib = "MiB"
	gib = "GiB"
)

func toMib(input string) float32 {
	results := sizePattern.FindStringSubmatch(input)
	value, err := strconv.ParseFloat(results[1], 32)
	if err != nil {
		log.Panic(err)
	}
	quantifier := results[2]

	switch quantifier {
	case mib:
		break
	case b:
		if value != 0 {
			value = value / (1024 * 1024)
		}
	case gib:
		value *= 1024
	default:
		log.Panicf("Unkown quantifier %s", quantifier)
	}

	return float32(value)
}
