package server

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"sync"

	"path/filepath"

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

type handle struct {
	// fileMap contains sometimes up to 200k entries
	fileMap map[string]record.Record
	sync.RWMutex
	length int
	// pattern is the name and type of the image pattern
	pattern       string
	virtualImages virtImgRoutes
}

func newHandle(ctx *cli.Context) *handle {
	csvFile := ctx.String("i")
	rec := record.GetCSVContent(csvFile)

	h := &handle{
		fileMap: make(map[string]record.Record),
		pattern: ctx.GlobalString("p"),
	}
	h.Lock()
	defer h.Unlock()

	if vif := ctx.String("imgconfig"); vif != "" {
		var err error
		if h.virtualImages, err = parseVirtualImageConfigFile(vif); err != nil {
			common.InfoErr("File %s contains error:\n%s\n", vif, err)
		}
	}

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

func (h *handle) root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, TextHTMLCharsetUTF8)
	fmt.Fprint(w, `<html>
	<head><title>Mediamock Index</title></head>
	<body>
		<a href="/json">JSON Index</a><br>
		<a href="/html">HTML Index</a><br>
		<a href="/debug/charts/">Debug Charts</a><br>
	</body>
	</html>`)
}

func (h *handle) rootJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	h.RLock()

	//if err := json.NewEncoder(w).Encode(h.virtualImages); err != nil {
	//	common.InfoErr("Failed to write JSON encoded virtual image configuration with error: %s\n", err)
	//}

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

func (h *handle) rootHTML(w http.ResponseWriter, r *http.Request) {
	h.RLock()
	defer h.RUnlock()
	w.Header().Set(ContentType, TextHTMLCharsetUTF8)
	fmt.Fprint(w, `<html>	<head><title>Mediamock Content Table</title></head>	<body><table>`)
	fmt.Fprint(w, `<thead><tr>
			<th>ModTime</th>
			<th nowrap>Width px</th>
			<th nowrap>Height px</th>
			<th>Link</th>
	</tr></thead><tbody>`)

	for path, vis := range h.virtualImages {
		for _, vi := range vis {
			_, err := fmt.Fprintf(w, `<tr>
			<td>n/a</td>
			<td>%d</td>
			<td>%d</td>
			<td><a href="%s" target="_blank">%s</a> Regex: %s</td>
		</tr>`,
				vi.Width,
				vi.Height,
				path, path, vi.regex,
			)
			if err != nil {
				common.InfoErr("Failed to write HTML table with error: %s\n", err)
			}

			if _, err := w.Write(brByte); err != nil {
				common.InfoErr("Failed to write brByte with error: %s\n", err)
			}
		}
	}

	var pathSlice = make(sort.StringSlice, len(h.fileMap))
	var i int
	for key, _ := range h.fileMap {
		pathSlice[i] = key
		i++
	}
	pathSlice.Sort()

	for _, key := range pathSlice {
		rec := h.fileMap[key]

		_, err := fmt.Fprintf(w, `<tr>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td><a href="%s" target="_blank">%s</a></td>
		</tr>`,
			rec.ModTime,
			rec.Width,
			rec.Height,
			rec.Path, rec.Path,
		)
		if err != nil {
			common.InfoErr("Failed to write HTML table with error: %s\n", err)
		}

		if _, err := w.Write(brByte); err != nil {
			common.InfoErr("Failed to write brByte with error: %s\n", err)
		}

	}

	fmt.Fprint(w, `</tbody></table></body>	</html>`)
}

// post or get request
// $ curl --data "file=media/catalog/product/1/2/120---microsoft-natural-ergonomic-keyboard-4000.jpg" http://127.0.0.1:4711/fileDetails
// and returns:
// {"Path":"media/catalog/product/1/2/120---microsoft-natural-ergonomic-keyboard-4000.jpg","ModTime":"2014-02-16T03:27:45+01:00","Width":5184,"Height":3456}
func (h *handle) fileDetails(w http.ResponseWriter, r *http.Request) {
	filePath := r.FormValue("file")
	if filePath == "" {
		http.NotFound(w, r)
		return
	}

	h.RLock()
	defer h.RUnlock()

	rec, ok := h.fileMap[filePath]
	if !ok {
		common.InfoErr("%s not found in CSV file\n", filePath)
		http.NotFound(w, r)
		return
	}
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	if err := rec.ToJSON(w); err != nil {
		common.InfoErr("Failed to write JSON with error: %s\n", err)
	}
}

func (h *handle) findRecord(w http.ResponseWriter, r *http.Request) (rec record.Record, matched bool) {
	h.RLock()
	defer h.RUnlock()

	// 1. try virtual image configuration from TOML
	// 2. try file map from scanned directory and the generated CSV file
	// 3. try to find a pattern in the URL which might look like an image size: 420x230
	// 4. print not found

	var path = r.URL.Path[1:]
	dir, file := filepath.Split(path)

	if vis, ok := h.virtualImages[dir]; ok {
		for _, vi := range vis {

			if vi.regex != nil {
				matched = vi.regex.MatchString(file)
				if matched {
					rec = record.NewRecordFields(h.pattern, path, vi.Width, vi.Height)
					return
				}

			} else {
				matched = true
				rec = record.NewRecordFields(h.pattern, path, vi.Width, vi.Height)
				return
			}
		}
	}

	rec, matched = h.fileMap[path]
	if !matched && len(path) > 0 && path[0] != '/' {
		rec, matched = h.fileMap["/"+path]
	}

	if false == matched {
		if false == common.IsImage(path) {
			http.NotFound(w, r)
			return
		}

		if width, height := common.FileSizeFromPath(path); width > 0 && height > 0 {
			matched = true
			rec = record.NewRecordFields(h.pattern, path, width, height)
		} else {
			http.NotFound(w, r)
			return
		}
	}

	return
}

func (h *handle) handler(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/":
		h.root(w, r)
		return
	case "/json":
		h.rootJSON(w, r)
		return
	case "/html":
		h.rootHTML(w, r)
		return
	case "/fileDetails":
		h.fileDetails(w, r)
		return
	case "/robots.txt":
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "User-agent: *\nDisallow: /")
		return
	}

	var path = r.URL.Path[1:]

	rec, ok := h.findRecord(w, r)
	if !ok {
		return
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
