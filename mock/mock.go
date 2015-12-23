package mock

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/record"
	"github.com/codegangsta/cli"
)

// @todo remove CSV/GZ and use boltDB

func ActionCLI(ctx *cli.Context) {
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

	if false == common.IsDir(targetDir) {
		common.UsageAndExit("Expecting a directory: %s", targetDir)
	}

	records := record.GetCSVContent(csvFile)

	var count = len(records)
	var recordChan = make(chan record.Record)
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(&wg, i, recordChan, targetDir)
	}

	var t = time.Now()
	var bar = common.InitPB(count)
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
