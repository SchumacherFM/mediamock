package main

import (
	"os"

	"time"

	"github.com/SchumacherFM/mediamock/analyze"
	"github.com/SchumacherFM/mediamock/common"
	"github.com/SchumacherFM/mediamock/server"
	"github.com/codegangsta/cli"
)

var (
	BUILD_DATE string
	VERSION    string // will be set via goxc from outside
	fileName   = func() (fn string) {
		var err error
		if fn, err = os.Hostname(); err == nil {
			fn = fn + "_"
		}
		return fn + "mediamock.csv.gz"
	}()
)

func main() {

	app := cli.NewApp()
	app.Name = "mediamock"
	if VERSION == "" {
		VERSION = "develop"
		BUILD_DATE = time.Now().String()
	}
	app.Version = VERSION + " by @SchumacherFM (compiled " + BUILD_DATE + ")"
	app.Usage = `reads your assets/media directory on your server and
               replicates it as a virtual structure on your development machine.
               On top can act as a proxy.

               $ mediamock help analyze|server|imgconfig for more options!
               `
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "q",
			Usage: "No output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "analyze",
			ShortName: "a",
			Usage: `Analyze the directory structure on you production server and write into a
                csv.gz file.`,
			Action: analyze.ActionCLI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "d",
					Value: ".",
					Usage: "Read this directory recursively and write into -o",
				},
				cli.StringFlag{
					Name:  "o",
					Value: common.TempDir() + fileName,
					Usage: "Write to this output file.",
				},
			},
		},
		{
			Name:      "server",
			ShortName: "s",
			Usage: `Server reads the csv.gz file and creates the assets/media structure on the fly
                 as a HTTP server. Does not write anything to your hard disk. Open URL / on the
                 server to retrieve a list of all files and folders.`,
			Action: server.ActionCLI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "img-config",
					Value: "",
					Usage: `Path to the configuration file for virtual image generation.

                 img-config defines a TOML configuration file which allows you to specify wilcard
                 image generation. You define a path to a directory and declare the image width
                 and height. All image http requests to that directory will have the same size.
                 Further more you can declare more occurences of the same directory and add a
                 regular expression to serve different width and height within that directory.
                 The image extension will be detected automatically. Type on the CLI:
                 '$ mediamock imgconfig' to see an example of a TOML config.
`,
				},
				cli.StringFlag{
					Name:  "img-pattern",
					Value: "icon",
					Usage: "Image pattern: happy, warm, rand, happytext, warmtext, a HTML hex value or icon",
				},
				cli.StringFlag{
					Name:  "url-prefix",
					Value: "",
					Usage: "Prefix in the URL path of the csv.gz file.",
				},
				cli.StringFlag{
					Name:  "csv",
					Value: "",
					Usage: "Source of csv.gz (file or URL)",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1:4711",
					Usage: "IP address or host name",
				},
				cli.StringFlag{
					Name:  "media-url",
					Value: "",
					Usage: `External URL to the base media directory.

                Apply this URL and mediamock will download the images and save them locally. If
                the remote image does not exists a mocked image will be generated. (Proxy Feature)
`,
				},
				cli.StringFlag{
					Name:  "media-cache",
					Value: "",
					Usage: `Local folder where to cache the downloaded images. (Proxy Feature)`,
				},
			},
		},
		{
			Name:  "imgconfig",
			Usage: `Prints an example TOML configuration file.`,
			Action: func(ctx *cli.Context) {
				println("A TOML configuration file might look like:")
				println(server.ExampleToml)
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
