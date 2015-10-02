package common

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"path/filepath"
	"strings"

	"github.com/mgutz/ansi"
)

func UsageAndExit(message string, args ...interface{}) {
	if message != "" {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, ansi.Color(message, "red"), args...)
		fmt.Fprintf(os.Stderr, "\n")
	}
	//flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func InfoErr(msg string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, ansi.Color(msg, "magenta"), args...); err != nil {
		panic(err)
	}
}

var fileSizePattern = regexp.MustCompile("([0-9]+)x([0-9]+)?")

func FileSizeFromPath(path string) (width, height int) {
	m := fileSizePattern.FindStringSubmatch(path)
	if len(m) == 3 {
		width, _ = strconv.Atoi(m[1])
		height, _ = strconv.Atoi(m[2])
	}
	if height == 0 {
		height = width
	}
	return
}

var ps = string(os.PathSeparator)

func ContainsFolderName(path string, names ...string) bool {
	for _, n := range names {
		if strings.Contains(path, ps+n+ps) {
			return true
		}
	}
	return false
}

func IsImage(path string) (ok bool) {

	switch filepath.Ext(path) {
	case ".png", ".gif", ".jpg", ".jpeg", ".ico":
		ok = true
	}
	return
}
