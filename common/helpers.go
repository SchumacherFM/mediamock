package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/rakyll/pb"
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

func Info(msg string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, ansi.Color(msg, "green"), args...); err != nil {
		panic(err)
	}
}

var fileSizePattern = regexp.MustCompile("([0-9]+)x([0-9]+)?")
var ps = string(os.PathSeparator)

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

// TempDir returns temporary directory ending with a path separatoe
func TempDir() string {
	dir := os.TempDir()
	if len(dir) > 0 && dir[len(dir)-1:] != ps {
		dir = dir + ps
	}
	return dir
}

var excludedFolders = []string{".svn", ".git"}

// ContainsFolderName checks if path contains a name. A name will be prepended
// with an OS specific path separator, e.g.: .svn becomes /.svn
func ContainsFolderName(path string, names ...string) bool {
	names = append(names, excludedFolders...)
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

// IsDir returns true if path is a directory
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	return fileInfo != nil && fileInfo.IsDir() && err == nil
}

// FileExists returns true if file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsHTTP checks if path starts with http or https
func IsHTTP(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// InitPB initializes a progress bar for the terminal
func InitPB(count int) *pb.ProgressBar {
	bar := pb.New(count)
	bar.ShowPercent = true
	bar.ShowBar = true
	bar.ShowCounters = true
	bar.ShowTimeLeft = true
	return bar
}
