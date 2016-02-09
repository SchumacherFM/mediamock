package server

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SchumacherFM/mediamock/common"
	"github.com/codegangsta/cli"
)

const maxHttpDownloadSize = 1024 * 1e3 * 6 // 6 MB

type proxy struct {
	url      string
	cacheDir string
}

func newProxy(ctx *cli.Context) *proxy {

	p := &proxy{
		url:      ctx.String("media-url"),
		cacheDir: ctx.String("media-cache"),
	}

	if p.url == "" {
		common.Info("Image proxying disabled. Media URL is empty.\n")
		return nil
	}

	if p.cacheDir == "" {
		h := fnv.New32a()
		if _, err := h.Write([]byte(p.url)); err != nil {
			panic(err)
		}
		p.cacheDir = fmt.Sprintf("%smediamock_proxy_%d%s", common.TempDir(), h.Sum32(), string(os.PathSeparator))
	}
	common.Info("Proxy started with remote URL %q and cache directory %q\n", p.url, p.cacheDir)

	return p
}

func (p *proxy) serveExistingFile(w http.ResponseWriter, r *http.Request) bool {
	reqFile := r.URL.Path[1:]
	cachedFile := p.cacheDir + reqFile

	if common.FileExists(cachedFile) {

		f, err := os.Open(cachedFile)
		if err != nil {
			if false == os.IsNotExist(err) {
				common.InfoErr("Cannot open %q because %s\n", cachedFile, err.Error())
			}
			return false
		}
		defer func() {
			if err := f.Close(); err != nil {
				common.InfoErr("Failed to close file %q with Error %s\n", reqFile, err)
			}
		}()

		addHeaders(filepath.Ext(reqFile), cacheUntil, w)
		if _, err := io.Copy(w, f); err != nil {
			common.InfoErr("Error copying file content to http response: %s with file %q", err, reqFile)
			return false
		}
		return true
	}
	return false
}

func (p *proxy) serveAndSaveRemoteFile(w http.ResponseWriter, r *http.Request) bool {
	reqFile := r.URL.Path[1:]
	cachedFile := p.cacheDir + reqFile
	remoteFile := p.url + reqFile

	resp, err := http.Get(remoteFile)
	defer func() {
		if resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				common.UsageAndExit("Failed to close URL %q with error: %s", remoteFile, err)
			}
		}
	}()

	if err != nil {
		common.InfoErr("Failed to download %q with error: %s\n", remoteFile, err)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		common.InfoErr("Failed to download %q with Code: %q\n", remoteFile, http.StatusText(resp.StatusCode))
		return false
	}

	if dcf := filepath.Dir(cachedFile); false == common.IsDir(dcf) {
		if err := os.MkdirAll(dcf, 0755); err != nil {
			common.InfoErr("Failed to create cache folder %q with Code: %s\n", dcf, err)
			return false
		}
	}

	fw, err := os.OpenFile(cachedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		common.InfoErr("Failed to open write file %q for URL %q with Error: %s\n", cachedFile, remoteFile, err)
		return false
	}
	defer func() {
		if err := fw.Close(); err != nil {
			common.InfoErr("Failed to close write file %q from URL %q with Error %s\n", cachedFile, remoteFile, err)
		}
	}()

	addHeaders(filepath.Ext(reqFile), cacheUntil, w)
	mw := io.MultiWriter(w, fw)
	if _, err := io.Copy(mw, io.LimitReader(resp.Body, maxHttpDownloadSize)); err != nil {
		common.InfoErr("Failed to write to http response and/or file to disk: %q with error: %s. Max Size: %d KBytes\n", remoteFile, err, maxHttpDownloadSize/1024)
		return false
	}

	return true
}

// pipe returns false == entry not found on server
func (p *proxy) serveProxy(w http.ResponseWriter, r *http.Request) bool {

	remoteFile := p.url + r.URL.Path[1:]
	if false == common.IsHTTP(remoteFile) {
		return false
	}

	if p.serveExistingFile(w, r) {
		return true
	}

	return p.serveAndSaveRemoteFile(w, r)

}
