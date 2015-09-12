package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const dirPerm os.FileMode = 0755

func mock(targetDir, csvFile string) {

	if targetDir != "" {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			if err := os.MkdirAll(targetDir, dirPerm); err != nil {
				usageAndExit("Failed to create directory %s with error: %s", targetDir, err)
			} else {
				fmt.Fprintf(os.Stdout, "Directory %s created\n", targetDir)
			}
		}

	}

	if targetDir == "" {
		targetDir = "."
	}

	if false == isDir(targetDir) {
		usageAndExit("Expecting a directory: %s", targetDir)
	}

	r := getCSVContent(csvFile)
	defer func() {
		if err := r.Close(); err != nil {
			usageAndExit("Failed to close file %s with error: %s", csvFile, err)
		}
	}()

	rz, err := gzip.NewReader(r)
	if err != nil {
		usageAndExit("Failed to create a GZIP reader from file %s with error: %s", csvFile, err)
	}
	defer func() {
		if err := rz.Close(); err != nil {
			usageAndExit("Failed to close file %s with error: %s", csvFile, err)
		}
	}()

	rc := csv.NewReader(rz)
	rc.Comma = ([]rune(csvSep))[0]

	var recordChan = make(chan record)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(&wg, i, recordChan, targetDir)
	}

	var i int
	var t = time.Now()
	for {

		raw, err := rc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read a record CSV data from file %s with error: %s\n", csvFile, err)
		}

		rec, err := newRecord(raw...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", csvFile, err)
		}

		recordChan <- rec

		if i%100 == 0 && i > 0 {
			fmt.Fprintf(os.Stdout, "%6d => %s\n", i, time.Now().Sub(t))
			t = time.Now()
		}
		i++
	}
	close(recordChan)
	wg.Wait()
	fmt.Fprintf(os.Stdout, "Created %d files\n", i+1)
}

func worker(wg *sync.WaitGroup, id int, rec <-chan record, targetDir string) {
	defer wg.Done()
	for { // or we could do for r := range rec { ... } what's better?
		r, ok := <-rec
		if !ok {
			return
		}
		if err := r.Create(targetDir); err != nil {
			fmt.Fprintf(os.Stderr, "Worker %d: Failed to create file: %s\n", id, err)
		}
		// fmt.Printf("Worker %d => %s\n", id, r.Path)
	}
}

func getCSVContent(csvFile string) io.ReadCloser {
	if isHTTP(csvFile) {
		resp, err := http.Get(csvFile)
		if err != nil {
			usageAndExit("Failed to download %s with error: %s", csvFile, err)
		}
		return resp.Body
	}

	fc, err := os.Open(csvFile)
	if err != nil {
		usageAndExit("Failed to open %s with error:%s", csvFile, err)
	}
	return fc
}

func isHTTP(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
