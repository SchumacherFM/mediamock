package server

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/codegangsta/cli"
	"io/ioutil"
)

const maxHttpDownloadSize = 1024 * 1e3 * 6 // 6 MB

type proxy struct {
	url      string
	cacheDir string
	mu       sync.Mutex
	cached   map[string]bool
}

func newProxy(ctx *cli.Context) *proxy {

	p := &proxy{
		url:      ctx.String("media-url"),
		cacheDir: ctx.String("media-cache"),
		cached:   make(map[string]bool),
	}

	if p.url == "" {
		common.InfoErr("Image proxying disabled. Media URL is empty.")
		return nil
	}

	if p.cacheDir == "" {
		h := fnv.New64a()
		if _, err := h.Write([]byte(p.url)); err != nil {
			panic(err)
		}
		p.cacheDir = fmt.Sprintf("%smediamock_proxy_%d%s", common.TempDir(), h.Sum64(), string(os.PathSeparator))
	}
	common.InfoErr("Proxy started from URL %q and cache directory %q", p.url, p.cacheDir)

	return p
}

// pipe returns false == entry not found on server
func (p *proxy) pipe(w http.ResponseWriter, r *http.Request) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	reqFile := r.URL.Path[1:]

	if _, ok := p.cached[reqFile]; ok {

		f, err := os.Open(p.cacheDir + reqFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		defer func() {
			if err := f.Close(); err != nil {
				common.InfoErr("Failed to close file %q with Error %s", reqFile,err)
			}
		}()

		addHeaders(filepath.Ext(reqFile), cacheUntil, w)
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, fmt.Sprintf("Error: %s with file %q", err, reqFile), http.StatusInternalServerError)
		}
		return true
	}

	remoteFile := p.url + reqFile
	if false == common.IsHTTP(remoteFile) {
		return false
	}

	resp, err := http.Get(remoteFile)
	defer func() {
		if resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				common.UsageAndExit("Failed to close URL %q with error: %s", remoteFile, err)
			}
		}
	}()

	if err != nil {
		common.InfoErr("Failed to download %q with error: %s", remoteFile, err)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		common.InfoErr("Failed to download %q with Code: %s", remoteFile, http.StatusText(resp.StatusCode))
		return false
	}

	fileData, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxHttpDownloadSize))

	fw, err := os.OpenFile(p.cacheDir+reqFile, os.O_WRONLY, 0600)
	if err !=nil {
		common.InfoErr("Failed to open file %q  for URL %q with Error: %s",p.cacheDir+reqFile, remoteFile, err)
		return false
	}
	defer func() {
		if err := fw.Close(); err != nil {
			common.InfoErr("Failed to close file %q with Error %s", reqFile,err)
		}
	}()

	addHeaders(filepath.Ext(reqFile), cacheUntil, w)
	mw := io.MultiWriter(w, fw)
	if _, err := mw.Write(fileData); err != nil {
		common.InfoErr("Failed to write to http response and/or file to disk: %q with error: %s", remoteFile, err)
		return false
	}

	return true
}
