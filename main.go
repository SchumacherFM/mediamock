package main

import (
	"os"

	"flag"
	"log"
	"strings"
)

const csvSep = "|"
const csvNewLine = "\n"

var (
	inFile  = flag.String("csv", "", "Read CSV data from this gzip file.")
	readDir = flag.String("dir", "", "Read this directory recursivly and write into -o")
	outFile = flag.String("o", "/tmp/mediamock.csv.gz", "Write to this file")
)

func main() {
	flag.Parse()

	if *inFile == "" && *readDir == "" {
		flag.Usage()
		os.Exit(0)
	}

	if *inFile != "" {
		mock(*inFile)
		return
	}

	if _, err := os.Stat(*readDir); os.IsNotExist(err) {
		log.Fatalf("No such file or directory: %s\n", *readDir)
	}
	analyze(*readDir, *outFile)

}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir() && err == nil
}

func isHTTP(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
