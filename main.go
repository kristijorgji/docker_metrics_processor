package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bye()
		os.Exit(0)
	}()

	log.Printf("Starting with Batch size %d, input path %s\n", batchSize, inputPath)

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

	log.Printf("[%s] Started processing\n", path)
	metrics := parser.Parse(path)
	log.Printf("[%s] Parsed, will insert into storage\n", path)

	var batch []*models.ServiceMetrics
	i := 0
	for ; i < len(metrics); i++ {
		batch = append(batch, &metrics[i])
		if i != 0 && i%batchSize == 0 {
			metricsRepository.InsertBatch(batch)
			log.Printf("[%s] Inserted %d rows", path, batchSize)
			batch = make([]*models.ServiceMetrics, 0)
		}
	}

	if len(batch) > 0 {
		metricsRepository.InsertBatch(batch)
		log.Printf("[%s] Inserted %d rows", path, batchSize)
		batch = nil
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
