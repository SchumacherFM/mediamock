package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

const fileName = "mediamock.csv.gz"

var (
	inFile  = flag.String("i", "", "")
	dir     = flag.String("d", "", "")
	outFile = flag.String("o", "", "")
	pattern = flag.String("p", "happy", "")
)

var usage = `Usage: mediamock options...

Options:
  -i  Read CSV data from this input URL/file.
  -d  Read this directory recursively and write into -o. If -i is provided
      generate all mocks in this directory. Default: current directory.
  -o  Write data into out file (optional, default a temp file).
  -p  Image pattern: happy (default), warm, rand or HTML hex value
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if flag.NFlag() < 1 {
		usageAndExit("")
	}

	if *inFile != "" {
		mock(*dir, *inFile)
		return
	}

	if *outFile == "" {
		*outFile = os.TempDir() + fileName
	}

	analyze(*dir, *outFile)

}

func usageAndExit(message string, args ...interface{}) {
	if message != "" {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, ansi.Color(message, "red"), args...)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func infoErr(msg string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, ansi.Color(msg, "magenta"), args...)
}
