package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// ContainsFolderName checks if path contains a name. A name will be prepended
// with an OS specific path separator, e.g.: .svn becomes /.svn
func ContainsFolderName(path string, names ...string) bool {
	for _, n := range names {
		if strings.Contains(path, ps+n) {
			return true
		}
	}
	return false
}

// IsImage checks if path to a file is an image by extracting the file extension
// and checking it against an internal list.
func IsImage(path string) (ok bool) {
	switch filepath.Ext(path) {
	case ".png", ".gif", ".jpg", ".jpeg", ".ico":
		ok = true
	}
	return
}
