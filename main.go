package main

import (
	"fmt"
	"log"
	"time"

	"./parser"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	start := time.Now()
	bye := func() {
		log.Printf("Execution took %s\n", time.Since(start))
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
		bye()
	}()

	metrics := parser.Parse("./input/2019-11-01.log")

	for _, element := range metrics {
		fmt.Println(element)
	}
}

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
}
