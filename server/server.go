package server

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/codegangsta/cli"
	_ "github.com/mkevac/debugcharts"
)

var (
	brByte     = []byte("\n")
	cacheUntil = time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)
)

func ActionCLI(ctx *cli.Context) {
	h := newHandle(ctx)
	fmt.Fprintf(os.Stdout, "Found %d entries in the CSV file\n", h.length)
	fmt.Fprintf(os.Stdout, "Server started: %s\n", ctx.String("host"))
	http.HandleFunc("/", h.handler)
	if err := http.ListenAndServe(ctx.String("host"), nil); err != nil {
		panic(err)
	}
}
