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

	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/record"
	"github.com/codegangsta/cli"
	"github.com/rakyll/pb"
)

// @todo remove CSV/GZ and use boltDB

func actionMock(ctx *cli.Context) {
	var targetDir, csvFile string

	targetDir = ctx.String("d")
	csvFile = ctx.String("i")

	if targetDir != "" {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			if err := os.MkdirAll(targetDir, record.DirPerm); err != nil {
				common.UsageAndExit("Failed to create directory %s with error: %s", targetDir, err)
			} else {
				fmt.Fprintf(os.Stdout, "Directory %s created\n", targetDir)
			}
		}
	}

	if targetDir == "" {
		targetDir = "."
	}

	if false == isDir(targetDir) {
		common.UsageAndExit("Expecting a directory: %s", targetDir)
	}

	records := getCSVContent(csvFile)

	var count = len(records)
	var recordChan = make(chan record.Record)
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(&wg, i, recordChan, targetDir)
	}

	var t = time.Now()
	var bar = initPB(count)
	bar.Start()

	for _, row := range records {
		rec, err := record.NewRecord(ctx.String("p"), row...)
		if err != nil {
			common.InfoErr("File %s contains error: %s\n", csvFile, err)
		}
		recordChan <- rec
		bar.Increment()
	}
	close(recordChan)
	wg.Wait()
	bar.Finish()
	fmt.Fprintf(os.Stdout, "Duration: %s\n", time.Now().Sub(t))
	fmt.Fprintf(os.Stdout, "You may run: find %s -type f -name \"*.jpg\" -exec jpegoptim {} + \n", targetDir)
}

func worker(wg *sync.WaitGroup, id int, rec <-chan record.Record, targetDir string) {
	defer wg.Done()
	for { // or we could do for r := range rec { ... } what's better?
		r, ok := <-rec
		if !ok {
			return
		}
		if err := r.CreateFile(targetDir); err != nil {
			common.InfoErr("Worker %d: Failed to create file: %s\n", id, err)
		}
		// fmt.Printf("Worker %d => %s\n", id, r.Path)
	}
}

func getCSVContent(csvFile string) [][]string {
	var rawRC io.ReadCloser
	if isHTTP(csvFile) {
		resp, err := http.Get(csvFile)
		if err != nil {
			common.UsageAndExit("Failed to download %s with error: %s", csvFile, err)
		}
		if resp.StatusCode != http.StatusOK {
			common.UsageAndExit("Server return non-200 status code: %s\nFailed to download %s", resp.Status, csvFile)
		}
		rawRC = resp.Body
	} else {
		var err error
		rawRC, err = os.Open(csvFile)
		if err != nil {
			common.UsageAndExit("Failed to open %s with error:%s", csvFile, err)
		}
	}
	defer func() {
		if err := rawRC.Close(); err != nil {
			common.UsageAndExit("Failed to close URL/file %s with error: %s", csvFile, err)
		}
	}()

	rz, err := gzip.NewReader(rawRC)
	if err != nil {
		common.UsageAndExit("Failed to create a GZIP reader from file %s with error: %s", csvFile, err)
	}
	defer func() {
		if err := rz.Close(); err != nil {
			common.UsageAndExit("Failed to close file %s with error: %s", csvFile, err)
		}
	}()

	rc := csv.NewReader(rz)
	rc.Comma = ([]rune(record.CSVSep))[0]

	records, err := rc.ReadAll()
	if err != nil {
		common.UsageAndExit("Failed to read CSV file %s with error: %s", csvFile, err)
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
	return bar
}
