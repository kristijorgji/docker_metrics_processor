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
		log.Printf("Total execution took %s\n", time.Since(start))
	}
	defer bye()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bye()
		os.Exit(0)
	}()

	log.Printf("Starting with batch size of %d, input path %s\n", batchSize, inputPath)

	metricsRepository = &repositories.MetricsRepository{}
	metricsRepository.Init()
	defer metricsRepository.Close()

	insertedRowsPerUnitOfWork := make(map[string]int)
	var wg sync.WaitGroup

	err := filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				wg.Add(1)
				go processLog(&wg, insertedRowsPerUnitOfWork, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()

	totalInsertedRows := 0
	for _, v := range insertedRowsPerUnitOfWork {
		totalInsertedRows += v
	}
	log.Printf("Inserted total %d rows", totalInsertedRows)
}

func processLog(wg *sync.WaitGroup, insertedRowsPerUnitOfWork map[string]int, path string) {
	defer timeTrack(time.Now(), fmt.Sprintf("[%s] Finished processing", path))
	defer wg.Done()

	log.Printf("[%s] Started processing\n", path)
	metrics := parser.Parse(path)
	log.Printf("[%s] Parsed, will insert into storage\n", path)

	var totalInsertedRows int = 0
	var batch []*models.ServiceMetrics
	i := 0
	for ; i < len(metrics); i++ {
		batch = append(batch, &metrics[i])
		batchLength := len(batch)
		if batchLength == batchSize {
			metricsRepository.InsertBatch(batch)
			log.Printf("[%s] Inserted %d rows", path, batchLength)
			totalInsertedRows += batchLength
			batch = make([]*models.ServiceMetrics, 0)
		}
	}

	batchLength := len(batch)
	if batchLength > 0 {
		metricsRepository.InsertBatch(batch)
		totalInsertedRows += batchLength
		log.Printf("[%s] Inserted %d rows", path, batchSize)
		batch = nil
	}

	log.Printf("[%s] Inserted total %d rows", path, totalInsertedRows)
	insertedRowsPerUnitOfWork[path] = totalInsertedRows

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
