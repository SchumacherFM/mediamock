package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

const csvSep = "|"
const csvNewLine = "\n"
const ps = string(os.PathSeparator)

type record struct {
	Path    string    // idx 0
	ModTime time.Time // idx 1
	Width   int       // idx 2
	Height  int       // idx 3
}

func newRecord(csv ...string) (record, error) {
	if len(csv) != 4 {
		return record{}, fmt.Errorf("Incorrect number %d of %d of CSV columns: %#v", len(csv), 4, csv)
	}

	w, _ := strconv.Atoi(csv[2])
	h, _ := strconv.Atoi(csv[3])
	t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", csv[1])

	return record{
		Path:    csv[0],
		ModTime: t,
		Width:   w,
		Height:  h,
	}, nil
}

func (r record) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s%s%s%s%d%s%d%s", r.Path, csvSep, r.ModTime, csvSep, r.Width, csvSep, r.Height, csvNewLine)
	return err
}

func (r record) Create(basePath string) error {

	d, f := r.getDirFile(basePath)
	if err := os.MkdirAll(d, dirPerm); err != nil {
		usageAndExit("Failed to create directory: %s with error: %s", d, err)
	}

	file, err := os.Create(d + f)
	if err != nil {
		return fmt.Errorf("Failed to create file %s%s", d, f)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file %s with error: %s\n", d+f, err)
		}
	}()

	switch filepath.Ext(f) {
	case ".png":
		png.Encode(file, r.generateImage())
	case ".jpg", ".jpeg":
		jpeg.Encode(file, r.generateImage(), &jpeg.Options{Quality: 1})
	case ".gif":
		gif.Encode(file, r.generateImage(), nil)
	default:
		if _, err := file.Write(nil); err != nil {
			return err
		}
	}

	return nil
}

func (r record) generateImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))

	var src image.Image
	switch *pattern {
	case "warm":
		src = &image.Uniform{colorful.WarmColor()}
	case "happy":
		src = &image.Uniform{colorful.HappyColor()}
	default:
		src = &image.Uniform{colorful.FastWarmColor()}
	}

	//	colorful.FastHappyColor()

	draw.Draw(img, img.Bounds(), src, image.ZP, draw.Src)
	return img
}

func (r record) getDirFile(base string) (dir, file string) {
	if false == strings.HasSuffix(base, ps) && false == strings.HasPrefix(r.Path, ps) {
		base = base + ps
	}

	dir, file = filepath.Split(base + r.Path)
	return
}
