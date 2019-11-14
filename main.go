package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"./models"
	"./parser"
	"./repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const defaultBatchSize = 1000

var batchSize int
var inputPath string

var metricsRepository *repositories.MetricsRepository

func main() {
	start := time.Now()
	setupEnv()
	bye := func() {
		log.Printf("Execution took %s\n", time.Since(start))
	}
	defer bye()

	metricsRepository = &repositories.MetricsRepository{}
	metricsRepository.Init()
	defer metricsRepository.Close()

	var wg sync.WaitGroup

	err := filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				wg.Add(1)
				go processLog(&wg, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
}

func processLog(wg *sync.WaitGroup, path string) {
	defer timeTrack(time.Now(), fmt.Sprintf("Finished processing %s\n", path))
	defer wg.Done()

	log.Printf("Started processing %s\n", path)
	metrics := parser.Parse(path)

	var batch []*models.ServiceMetrics
	for i := 0; i < len(metrics); i++ {
		batch = append(batch, &metrics[i])
		if i != 0 && i%batchSize == 0 {
			metricsRepository.InsertBatch(batch)
			batch = make([]*models.ServiceMetrics, 0)
		}
	}
}

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	batchSizeFlag := flag.Int("batchSize", defaultBatchSize, "the number of metrics to insert in storage in one batch")
	flag.StringVar(&inputPath, "inputPath", "input/", "folder path where logs are located")
	flag.Parse()

	batchSize = *batchSizeFlag
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s. Took %s", name, elapsed)
}
