package analyze

import (
	"compress/gzip"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/record"
	"github.com/rakyll/pb"
)

type traverse struct {
	cliWriter  io.Writer
	outfile    string
	basePath   string
	outF       io.WriteCloser
	outW       io.WriteCloser
	bar        *pb.ProgressBar
	workerRec  chan record.Record
	workerStop chan struct{}
}

func newTraverse(cliWriter io.Writer, path, outfile string, barMaxCount int) *traverse {
	w := &traverse{
		cliWriter:  cliWriter,
		outfile:    outfile,
		basePath:   path,
		workerRec:  make(chan record.Record),
		workerStop: make(chan struct{}),
		bar:        common.InitPB(barMaxCount),
	}

	var err error
	w.outF, err = os.Create(outfile)
	if err != nil {
		common.UsageAndExit("Failed to open %s with error: %s", outfile, err)
	}
	w.outW = gzip.NewWriter(w.outF)
	// using w.bar.Output with a io.Writer causes some funny print outs.
	w.bar.NotPrint = w.cliWriter == ioutil.Discard
	w.bar.Start()
	go w.workerWriter()
	return w
}

func (w *traverse) workerWriter() {
	for {
		select {
		case rec, ok := <-w.workerRec:
			if !ok {
				return
			}

			if err := rec.Write(w.outW); err != nil {
				fmt.Fprintf(w.cliWriter, "Error when writing to file %s: %s\n", w.outfile, err)
			}

			w.bar.Increment()
		case <-w.workerStop:
			return
		}
	}
}

func (w *traverse) close() {

	w.workerStop <- struct{}{}
	close(w.workerRec)

	if err := w.outW.Close(); err != nil {
		common.InfoErr("GZIP close error: %s\n", err)
	}
	if err := w.outF.Close(); err != nil {
		common.InfoErr("File close error: %s\n", err)
	}

	w.bar.Finish()
	fmt.Fprintf(w.cliWriter, "Wrote to file: %s\n", w.outfile)
}

func (w *traverse) walkFn(path string, info os.FileInfo, err error) error {
	rel := getRelative(w.basePath, path)

	if rel == "" || info.IsDir() {
		return nil
	}

	if common.ContainsFolderName(rel) {
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

	w.workerRec <- record.Record{
		Path:    rel,
		ModTime: info.ModTime(),
		Width:   imgWidth,
		Height:  imgHeight,
	}
	return nil
}

func getRelative(basePath, path string) string {
	path = filepath.Clean(path)
	if basePath == "" {
		return path
	}
	parts := strings.Split(path, basePath)

	if len(parts) < 2 {
		return ""
	}
	return parts[1]
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
