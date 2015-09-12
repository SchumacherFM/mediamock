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

	"github.com/cheggaaa/pb"
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

	records := getCSVContent(csvFile)

	var count = len(records)
	var recordChan = make(chan record)
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(&wg, i, recordChan, targetDir)
	}

	var t = time.Now()
	var bar = initPB(count)

	for _, row := range records {
		rec, err := newRecord(row...)
		if err != nil {
			infoErr("%s\n", csvFile, err)
		}
		recordChan <- rec
		bar.Increment()
	}
	close(recordChan)
	wg.Wait()
	bar.Finish()
	fmt.Fprintf(os.Stdout, "Duration: %s\n", time.Now().Sub(t))
}

func worker(wg *sync.WaitGroup, id int, rec <-chan record, targetDir string) {
	defer wg.Done()
	for { // or we could do for r := range rec { ... } what's better?
		r, ok := <-rec
		if !ok {
			return
		}
		if err := r.Create(targetDir); err != nil {
			infoErr("Worker %d: Failed to create file: %s\n", id, err)
		}
		// fmt.Printf("Worker %d => %s\n", id, r.Path)
	}
}

func getCSVContent(csvFile string) [][]string {
	var rawRC io.ReadCloser
	if isHTTP(csvFile) {
		resp, err := http.Get(csvFile)
		if err != nil {
			usageAndExit("Failed to download %s with error: %s", csvFile, err)
		}
		if resp.StatusCode != http.StatusOK {
			usageAndExit("Server return non-200 status code: %s\nFailed to download %s", resp.Status, csvFile)
		}
		rawRC = resp.Body
	} else {
		var err error
		rawRC, err = os.Open(csvFile)
		if err != nil {
			usageAndExit("Failed to open %s with error:%s", csvFile, err)
		}
	}
	defer func() {
		if err := rawRC.Close(); err != nil {
			usageAndExit("Failed to close URL/file %s with error: %s", csvFile, err)
		}
	}()

	rz, err := gzip.NewReader(rawRC)
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

	records, err := rc.ReadAll()
	if err != nil {
		usageAndExit("Failed to read CSV file %s with error: %s", csvFile, err)
	}

	return records
}

func isHTTP(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func initPB(count int) *pb.ProgressBar {
	bar := pb.New(count)
	bar.ShowPercent = true
	bar.ShowBar = true
	bar.ShowCounters = true
	bar.ShowTimeLeft = true
	return bar.Start()
}
