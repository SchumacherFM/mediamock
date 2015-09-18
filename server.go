package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/record"
	"github.com/codegangsta/cli"
	_ "github.com/mkevac/debugcharts"
)

const (
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
	ContentType     = "Content-Type"

	CharsetUTF8                = "charset=utf-8"
	ApplicationJSON            = "application/json"
	ApplicationJSONCharsetUTF8 = ApplicationJSON + "; " + CharsetUTF8

	TextHTML             = "text/html"
	TextHTMLCharsetUTF8  = TextHTML + "; " + CharsetUTF8
	TextPlain            = "text/plain"
	TextPlainCharsetUTF8 = TextPlain + "; " + CharsetUTF8
)

var (
	brByte     = []byte("\n")
	cacheUntil = time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)
)

func actionServer(ctx *cli.Context) {
	h := newHandle(ctx)
	fmt.Fprintf(os.Stdout, "Found %d entries in the CSV file\n", h.length)
	fmt.Fprintf(os.Stdout, "Server started: %s\n", ctx.String("host"))
	http.HandleFunc("/", h.handler)
	if err := http.ListenAndServe(ctx.String("host"), nil); err != nil {
		panic(err)
	}
}

type handle struct {
	// fileMap contains sometimes up to 200k entries
	fileMap map[string]record.Record
	sync.RWMutex
	length  int
	pattern string
}

func newHandle(ctx *cli.Context) *handle {
	csvFile := ctx.String("i")
	rec := getCSVContent(csvFile)

	h := &handle{
		fileMap: make(map[string]record.Record),
		pattern: ctx.GlobalString("p"),
	}
	h.Lock()
	defer h.Unlock()

	h.fileMap["favicon.ico"] = record.NewRecordFields("happy", "favicon.ico", 16, 16)
	for _, row := range rec {
		rec, err := record.NewRecord(h.pattern, row...)
		if err != nil {
			common.InfoErr("File %s contains error: %s\n", csvFile, err)
		}
		rec.Path = ctx.String("urlPrefix") + rec.Path
		h.fileMap[rec.Path] = rec
	}
	h.length = len(rec)
	return h
}

// root generates a JSON stream of all files
func (h *handle) root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	h.RLock()
	for _, rec := range h.fileMap {
		if err := rec.ToJSON(w); err != nil {
			common.InfoErr("Failed to write JSON with error: %s\n", err)
		}
		if _, err := w.Write(brByte); err != nil {
			common.InfoErr("Failed to write JSON with error: %s\n", err)
		}
	}
	h.RUnlock()
}

func (h *handle) handler(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/":
		h.root(w, r)
		return
	case "/robots.txt":
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "User-agent: *\nDisallow: /")
		return
	}

	h.RLock()
	defer h.RUnlock()

	var path = r.URL.Path[1:]
	rec, ok := h.fileMap[path]
	if !ok {

		if false == common.IsImage(path) {
			http.NotFound(w, r)
			return
		}

		width, height := common.FileSizeFromPath(path)
		if width > 0 && height > 0 {
			rec = record.NewRecordFields(h.pattern, path, width, height)
		}

	}

	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", rec.ModTime.Format(http.TimeFormat))
	w.Header().Set("Expires", cacheUntil)

	switch rec.FileExt() {
	case ".gif":
		w.Header().Set(ContentType, "image/gif")
	case ".png":
		w.Header().Set(ContentType, "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set(ContentType, "image/jpeg")
	case ".css":
		w.Header().Set(ContentType, "text/css")
	case ".js":
		w.Header().Set(ContentType, "application/javascript")
	case ".pdf":
		w.Header().Set(ContentType, "application/pdf")
	case ".txt":
		w.Header().Set(ContentType, TextPlain)
	case ".html":
		w.Header().Set(ContentType, TextHTML)
	}

	rec.CreateContent(path, w)
}
