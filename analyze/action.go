package analyze

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

	w := newTraverse(path, outfile)
	defer w.close()
	if err := walk.Walk(path, w.walkFn); err != nil {
		common.UsageAndExit("Walk Error: %s", err)
	}
	fmt.Fprintf(os.Stdout, "Wrote to file: %s\n", outfile)

}
