package main

import (
	"os"
	"strings"
)

func mock(csvFile, dir string) {
	if dir == "" {
		dir = "."
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		usageAndExit("No such file or directory: %s\n", *dir)
	}

	if false == isDir(dir) {
		usageAndExit("Expecting a directory: %s", dir)
	}

}



func isHTTP(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
