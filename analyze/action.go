package analyze

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"io"
	"io/ioutil"
	"os"

	"github.com/MichaelTJones/walk"
	"github.com/SchumacherFM/mediamock/common"
	"github.com/codegangsta/cli"
)

func ActionCLI(ctx *cli.Context) {

	var path, outfile string
	path = ctx.String("d")
	outfile = ctx.String("o")

	if false == common.IsDir(path) {
		common.UsageAndExit("Expecting an existing directory: %s", path)
	}

	var cliWriter io.Writer
	cliWriter = os.Stdout
	if ctx.GlobalBool("q") {
		cliWriter = ioutil.Discard
	}

	wc := newWalkCount(path)
	if err := walk.Walk(path, wc.walkFn); err != nil {
		common.UsageAndExit("Walk Counter Error: %s", err)
	}

	if path == "." {
		path = ""
	}

	w := newTraverse(cliWriter, path, outfile, wc.fileCount)
	defer w.close()
	if err := walk.Walk(path, w.walkFn); err != nil {
		common.UsageAndExit("Walk Error: %s", err)
	}
}
