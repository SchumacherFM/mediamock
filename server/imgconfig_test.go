package server

import (
	"go/build"
	"path/filepath"
	"testing"
)

func TestParseImgConfigFile(t *testing.T) {
	path := filepath.Join(build.Default.GOPATH, "src", "github.com", "SchumacherFM", "mediamock", "server", "testdata", "imgconfig.toml")

	tc, err := parseImgConfigFile(path)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n%#v\n", tc)
}
