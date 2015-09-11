package main

import (
	"compress/gzip"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func analyze(path, outfile string) {
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		usageAndExit("No such file or directory: %s\n", *dir)
	}

	if false == isDir(path) {
		usageAndExit("Expecting a directory: %s", path)
	}

	w := newWalk(path, outfile)
	defer w.close()
	if err := filepath.Walk(path, w.walkFn); err != nil {
		usageAndExit("Walk Error: %s", err)
	}
	fmt.Fprintf(os.Stdout, "Wrote to file: %s\n", outfile)

}

type walk struct {
	basePath string
	outF     io.WriteCloser
	outW     io.WriteCloser
}

func newWalk(path, outfile string) *walk {

	w := &walk{
		basePath: path,
	}

	var err error
	w.outF, err = os.Create(outfile)
	if err != nil {
		usageAndExit("Failed to open %s with error: %s", outfile, err)
	}
	w.outW = gzip.NewWriter(w.outF)

	return w
}

func (w *walk) close() {
	if err := w.outW.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "GZIP close error: %s\n", err)
	}
	if err := w.outF.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "File close error: %s\n", err)
	}
}

func (w *walk) getRelative(path string) string {
	path = filepath.Clean(path)
	parts := strings.Split(path, w.basePath)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func (w *walk) walkFn(path string, info os.FileInfo, err error) error {
	rel := w.getRelative(path)
	if rel == "" || info.IsDir() {
		return nil
	}

	var imgWidth, imgHeight int
	ext := strings.ToLower(filepath.Ext(rel))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		imgWidth, imgHeight = getImageDimension(path)
		//	default:
		//		log.Println(rel, ext)
	}

	return w.write(rel, info.ModTime(), imgWidth, imgHeight)
}

func (w *walk) write(path string, mod time.Time, width, height int) error {
	_, err := fmt.Fprintf(w.outW, "%s%s%s%s%d%s%d%s", path, csvSep, mod, csvSep, width, csvSep, height, csvNewLine)
	return err
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open image: %s\n", err)
		return 0, 0

	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Image %s decoding error: %s\n", imagePath, err)
		return 0, 0
	}

	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Close error: %s: %s\n", imagePath, err)
		return 0, 0
	}
	return image.Width, image.Height
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir() && err == nil
}
