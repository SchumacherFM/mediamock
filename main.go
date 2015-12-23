package main

import (
	"os"

	"github.com/SchumacherFM/mediamock/analyze"
	"github.com/SchumacherFM/mediamock/mock"
	"github.com/SchumacherFM/mediamock/server"
	"github.com/codegangsta/cli"
)

const fileName = "mediamock.csv.gz"

var (
	Version = "v0.3.0"
)

func main() {

	app := cli.NewApp()
	app.Name = "mediamock"
	app.Version = Version + " by @SchumacherFM"
	app.Usage = `reads your assets/media directory on your server and
               replicates that structure on your development machine.

               $ mediamock help analyze|mock|server for more options!
               `
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p",
			Value: "happy",
			Usage: "Image pattern: happy, warm, rand, happytext, warmtext HTML hex value or icon",
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
					Value: os.TempDir() + fileName,
					Usage: "Write to this output file.",
				},
			},
		},
		{
			Name:      "mock",
			ShortName: "m",
			Usage: `Mock reads the csv.gz file and recreates the files and folders. If a file represents
	an image, it will be created with a tiny file size and correct width x height.`,
			Action: mock.ActionCLI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i",
					Value: "",
					Usage: "Read csv.gz from this input file or input URL.",
				},
				cli.StringFlag{
					Name:  "d",
					Value: ".",
					Usage: "Create structure in this directory.",
				},
			},
		},
		{
			Name:      "server",
			ShortName: "s",
			Usage: `Server reads the csv.gz file and creates the assets/media structure on the fly
	as a HTTP server. Does not write anything to your hard disk. Open URL / on the server
	to retrieve a list of all files and folders.`,
			Action: server.ActionCLI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "urlPrefix",
					Value: "",
					Usage: "Prefix in the URL path",
				},
				cli.StringFlag{
					Name:  "i",
					Value: "",
					Usage: "Read csv.gz from this input file or input URL.",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1:4711",
					Usage: "IP address or host name",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
