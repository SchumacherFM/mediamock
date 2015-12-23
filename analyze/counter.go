package analyze

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"sync"

	"github.com/SchumacherFM/mediamock/common"
)

type walkCount struct {
	basePath string

	mu        sync.Mutex
	fileCount int
}

func newWalkCount(path string) *walkCount {
	return &walkCount{
		basePath: path,
	}
}

func (w *walkCount) walkFn(path string, info os.FileInfo, err error) error {
	rel := getRelative(w.basePath, path)

	if rel == "" || info.IsDir() || common.ContainsFolderName(rel) {
		return nil
	}

	w.mu.Lock()
	w.fileCount++
	w.mu.Unlock()

	return nil
}
