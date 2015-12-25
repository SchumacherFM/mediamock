package server

import (
	"go/build"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImgConfigFile(t *testing.T) {
	path := filepath.Join(build.Default.GOPATH, "src", "github.com", "SchumacherFM", "mediamock", "server", "testdata", "imgconfig.toml")

	tc, err := parseVirtualImageConfigFile(path)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tc, 2)
	assert.NotNil(t, tc["media/captcha/admin/"][0].regex)
	assert.NotNil(t, tc["media/captcha/admin/"][1].regex)
	assert.Exactly(t, 250, tc["media/captcha/base/"][0].Width)
}

func TestParseImgConfigFileError(t *testing.T) {
	path := filepath.Join(build.Default.GOPATH, "src", "github.com", "SchumacherFM", "mediamock", "server", "testdata", "imgconfigErr.toml")
	tc, err := parseVirtualImageConfigFile(path)
	assert.EqualError(t, err, "Entry\n'struct { Name string; Regex string; Width int; Height int }{Name:\"media/captcha/error/\", Regex:\"x(a-z]+\", Width:350, Height:120}'\ncannot be compiled to a valid regular expression\nError: error parsing regexp: missing closing ): `x(a-z]+`")
	assert.Nil(t, tc)
}
