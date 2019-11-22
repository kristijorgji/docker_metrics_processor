// Package repositories @author kristi.jorgji@flaconi.de created on 08.11.19
package repositories

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"../models"
)

// MetricsRepository fetch customers from data source
type MetricsRepository struct {
	db *sql.DB
}

// Init create an instance of repository
func (repository *MetricsRepository) Init() {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Panic(err.Error())
	}

	repository.db = db
}

// Close closes the connections
func (repository *MetricsRepository) Close() {
	err := repository.db.Close()
	if err != nil {
		log.Panic(err.Error())
	}
}

// InsertBatch into mysql
func (repository *MetricsRepository) InsertBatch(metrics []*models.ServiceMetrics) {
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO services VALUES")

	beforeLastIndex := len(metrics) - 1
	for i := 0; i < len(metrics); i++ {
		metric := *metrics[i]
		buffer.WriteString(
			fmt.Sprintf(
				"('%s', '%s', '%s', %f, %f, %f, %f)",
				metric.Datetime,
				metric.ContainerID,
				metric.ContainerName,
				metric.CPUPercentage,
				metric.MemoryUsageInMib,
				metric.MemoryLimitInMib,
				metric.MemoryPercentage,
			),
		)
		if i < beforeLastIndex {
			buffer.WriteString(",")
		}
	}

	insert, err := repository.db.Query(buffer.String())
	if err != nil {
		log.Panic(err.Error())
	}
	insert.Close()
}
