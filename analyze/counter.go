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
	mu        sync.Mutex
	dirCount  int64
	fileCount int64
}

func (w *walkCount) walkFn(path string, info os.FileInfo, err error) error {

	if common.ContainsFolderName(info.Name()) {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if info.IsDir() {
		w.dirCount++
		return nil
	}

	w.fileCount++

	return nil
}
