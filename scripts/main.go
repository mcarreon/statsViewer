package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./Close Evasive 1 - Challenge - 2025.01.17-03.50.56 Stats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Process the records
	for _, record := range records {
		fmt.Println(record) // Each record is a slice of strings
	}

	fmt.Println("done")
}
