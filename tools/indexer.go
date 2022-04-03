package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var (
	indexName, filename, urls string

	numWorkers int
	flushBytes int
	numItems   int
)

func init() {

	flag.StringVar(&urls, "urls", "http://localhost:9200", "filename name")
	flag.StringVar(&filename, "filename", "", "filename name")
	flag.StringVar(&indexName, "index", "sg1_index", "Index name")
	flag.IntVar(&numWorkers, "workers", runtime.NumCPU(), "Number of indexer workers")
	flag.IntVar(&flushBytes, "flush", 5e+6, "Flush threshold in bytes")
	flag.IntVar(&numItems, "count", 10000, "Number of documents to generate")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func readFileLines(logfile string, fnCbk func(string)) {
	f, err := os.OpenFile(logfile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("read file line error: %v", err)
			return
		}

		go fnCbk(line[:len(line)-1])
	}
}

// ./indexer -filename="/Users/51pwn/sgk1/BreachCompilation/data/0/0"
// ls /Users/51pwn/sgk1/BreachCompilation/data/0/?|xargs -I % ./indexer -filename="%"
// find  /Users/51pwn/sgk1/BreachCompilation/data|xargs -I % ./indexer -filename="%"
func main() {
	log.SetFlags(0)
	var wg sync.WaitGroup
	var (
		countSuccessful uint64
		err             error
	)

	log.Printf(
		"\x1b[1mBulkIndexer\x1b[0m: documents [%s] workers [%d] flush [%s]",
		humanize.Comma(int64(numItems)), numWorkers, humanize.Bytes(uint64(flushBytes)))
	log.Println(strings.Repeat("▁", 65))
	// Use a third-party package for implementing the backoff function
	retryBackoff := backoff.NewExponentialBackOff()

	// Create the Elasticsearch client
	// NOTE: For optimal performance, consider using a third-party HTTP transport package.
	//       See an example in the "benchmarks" folder.
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			urls,
		},
		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},

		// Configure the backoff function
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		MaxRetries: 5,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// Create the BulkIndexer
	// NOTE: For optimal performance, consider using a third-party JSON decoding package.
	//       See an example in the "benchmarks" folder.
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    numWorkers,       // The number of worker goroutines
		FlushBytes:    int(flushBytes),  // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}
	start := time.Now().UTC()
	readFileLines(filename, func(line1 string) {
		go func() {
			wg.Add(1)
			defer wg.Done()
			line := line1
			reg1 := regexp.MustCompile(`(^[^:;]+)`)
			var (
				uid, pswd string
			)
			uid = ""
			if reg1 != nil {
				result1 := reg1.FindAllStringSubmatch(line, -1)
				for _, text := range result1 {
					uid = text[1]
					break
				}
			}
			if "" == uid {
				n := strings.Index(line, ":")
				i := strings.Index(line, ";")
				x := len(line) - 1
				if 1 < n {
					x = n
				}
				if 1 < i && x > i {
					x = i
				}

				uid = line[:x]
			}
			pswd = line[len(uid)+1:]
			data := map[string]interface{}{"uid": uid, "pswd": []string{pswd}}
			bData, err := json.Marshal(data)
			if err != nil {
				return
			}
			fmt.Println(string(bData))
			if "" != uid {
				err = bi.Add(
					context.Background(),
					esutil.BulkIndexerItem{
						// Action field configures the operation to perform (index, create, delete, update)
						Action: "index",
						// DocumentID is the (optional) document ID
						DocumentID: uid,
						// Body is an `io.Reader` with the payload
						Body: bytes.NewReader(bData),
						// OnSuccess is called for each successful operation
						OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
							atomic.AddUint64(&countSuccessful, 1)
						},
						// OnFailure is called for each failed operation
						OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
							if err != nil {
								log.Printf("ERROR: %s", err)
							} else {
								log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
							}
						},
					},
				)
				if err != nil {
					log.Fatalf("Unexpected error: %s", err)
				}
			}
		}()
	})
	wg.Wait()
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}
