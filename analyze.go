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

	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/record"
	"github.com/codegangsta/cli"
)

func actionAnalyze(ctx *cli.Context) {

	var path, outfile string
	path = ctx.String("d")
	outfile = ctx.String("o")

	if false == isDir(path) {
		common.UsageAndExit("Expecting an existing directory: %s", path)
	}

	w := newWalk(path, outfile)
	defer w.close()
	if err := filepath.Walk(path, w.walkFn); err != nil {
		common.UsageAndExit("Walk Error: %s", err)
	}
	fmt.Fprintf(os.Stdout, "Wrote to file: %s\n", outfile)

}

type walk struct {
	basePath string
	outF     io.WriteCloser
	outW     io.WriteCloser
}

func newWalk(path, outfile string) *walk {
	if path == "." {
		path = ""
	}
	w := &walk{
		basePath: path,
	}

	var err error
	w.outF, err = os.Create(outfile)
	if err != nil {
		common.UsageAndExit("Failed to open %s with error: %s", outfile, err)
	}
	w.outW = gzip.NewWriter(w.outF)

	return w
}

func (w *walk) close() {
	if err := w.outW.Close(); err != nil {
		common.InfoErr("GZIP close error: %s\n", err)
	}
	if err := w.outF.Close(); err != nil {
		common.InfoErr("File close error: %s\n", err)
	}
}

func (w *walk) getRelative(path string) string {
	path = filepath.Clean(path)
	if w.basePath == "" {
		return path
	}
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

	if common.ContainsFolderName(rel, ".svn", ".git") {
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

	r := record.Record{
		Path:    rel,
		ModTime: info.ModTime(),
		Width:   imgWidth,
		Height:  imgHeight,
	}
	return r.Write(w.outW)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		common.InfoErr("Cannot open image: %s\n", err)
		return 0, 0

	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		common.InfoErr("Image %s decoding error: %s\n", imagePath, err)
		return 0, 0
	}

	if err := file.Close(); err != nil {
		common.InfoErr("Close error: %s: %s\n", imagePath, err)
		return 0, 0
	}
	return image.Width, image.Height
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir() && err == nil
}
