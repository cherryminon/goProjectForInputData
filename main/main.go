package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Record struct {
	Phone string `json:"phone"`
	Uid   string `json:"uid"`
}

func main() {
	esURL := "http://127.0.0.1:18074" // Replace with your Elasticsearch instance URL
	username := "elastic"             // Replace with your Elasticsearch username
	password := "LasNoches7"          // Replace with your Elasticsearch password
	indexNum := 0
	cfg := elasticsearch.Config{
		Addresses: []string{
			esURL,
		},
		Username: username,
		Password: password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	file, err := os.Open("/data/relation/wb.txt") // Replace with your txt file path
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup

	// Create a context object for the goroutines
	ctx := context.Background()

	// Bulk indexation
	var recordsBatch []string
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "\t")
		record := &Record{
			Phone: s[0],
			Uid:   s[1],
		}

		data, _ := json.Marshal(record)
		recordsBatch = append(recordsBatch, string(data))

		// If we have accumulated 1000 records, send them to Elasticsearch
		if len(recordsBatch) == 80000 {
			wg.Add(1)
			go func(batch []string) {
				defer wg.Done()
				BulkIndex(ctx, es, "wb_phone", batch)
				indexNum += 80000
				// Introduce some delay to prevent "Too Many Requests" error
				fmt.Println("本次处理的条数为", indexNum, "进入30秒休眠")
				time.Sleep(30 * time.Millisecond)
			}(recordsBatch)
			// Reset the batch
			recordsBatch = []string{}
		}
	}

	// Index remaining records
	if len(recordsBatch) > 0 {
		BulkIndex(ctx, es, "wb_phone", recordsBatch)
	}

	wg.Wait()

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}

func BulkIndex(ctx context.Context, es *elasticsearch.Client, index string, records []string) {
	// Prepare the data for the Bulk API
	var buf strings.Builder
	for _, record := range records {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
			},
		}

		// Append meta data
		metaJson, _ := json.Marshal(meta)
		buf.WriteString(string(metaJson) + "\n")

		// Append record data
		buf.WriteString(record + "\n")
	}

	// Send the request to Elasticsearch
	res, err := es.Bulk(strings.NewReader(buf.String()), es.Bulk.WithContext(ctx))
	if err != nil {
		log.Fatalf("Failure indexing batch: %s", err)
	}
	if res.IsError() {
		log.Printf("Failure indexing batch: %s", res.String())
	}
}
