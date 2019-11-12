package parser

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"
)

// ServiceMetrics Contains metrics about a docker service and the time
type ServiceMetrics struct {
	Datetime         string
	ContainerID      string
	ContainerName    string
	CPUPercentage    string
	MemoryUsage      string
	MemoryLimit      string
	MemoryPercentage string
}

const mysqlDateFormat = "2006-01-02 15:04:05.000000"

// Parse the file to the proper data structure
func Parse(filePath string) []ServiceMetrics {
	readFile, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var dateTimePattern = regexp.MustCompile("(\\d{4}-\\d{2}-\\d{2}.{15})")
	var dataPattern = regexp.MustCompile("(?P<containerId>.{12}) {2,}(?P<containerName>[^ ]+) {2,}(?P<cpuPercentage>[^ ]+%) {2,}(?P<memoryUsage>[^ ]+) / (?P<memoryLimit>[^ ]+) {2,}(?P<memoryPercentage>[^ ]+%) {2,}")

	var metrics []ServiceMetrics
	var serviceMetrics *ServiceMetrics = new(ServiceMetrics)

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
			serviceMetrics.CPUPercentage = result["cpuPercentage"]
			serviceMetrics.MemoryUsage = result["memoryUsage"]
			serviceMetrics.MemoryLimit = result["memoryLimit"]
			serviceMetrics.MemoryPercentage = result["memoryPercentage"]

			metrics = append(metrics, *serviceMetrics)
			serviceMetrics = new(ServiceMetrics)
		}

		i++
	}

	readFile.Close()

	return metrics
}
