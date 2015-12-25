package server

import (
	"regexp"

	"fmt"

	"github.com/BurntSushi/toml"
)

type virtImgRoutes map[string][]virtImgRoute

type virtImgRoute struct {
	regex  *regexp.Regexp
	Width  int
	Height int
}

func parseVirtualImageConfigFile(path string) (virtImgRoutes, error) {

	var cfg = struct {
		Dirs []struct {
			Name   string
			Regex  string
			Width  int
			Height int
		}
	}{}

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	r := make(virtImgRoutes)

	for _, c := range cfg.Dirs {

		if c.Width < 1 || c.Height < 1 || len(c.Name) < 1 {
			fmt.Printf("Virtual image entry\n'%#v'\nignored because of an invalid directory name, width or height.\n", c)
			continue
		}

		if c.Name[len(c.Name)-1:] != "/" {
			c.Name = c.Name + "/"
		}

		if _, ok := r[c.Name]; !ok {
			r[c.Name] = make([]virtImgRoute, 0, 1)
		}

		var reg *regexp.Regexp
		if c.Regex != "" {
			var err error
			reg, err = regexp.Compile(c.Regex)
			if err != nil {
				return nil, fmt.Errorf("Entry\n'%#v'\ncannot be compiled to a valid regular expression\nError: %s", c, err)
			}
		}
		r[c.Name] = append(r[c.Name], virtImgRoute{
			regex:  reg,
			Width:  c.Width,
			Height: c.Height,
		})

	}

	return r, nil
}

const ExampleToml = `[[dirs]]
name = "media/directory1/base"
width = 250
height = 120

[[dirs]]
name = "media/directory1/admin"
# regex matches the full file name
regex = "x[a-z]+"
width =  350
height = 120

[[dirs]]
name = "media/directory1/admin"
# regex matches the full file name
regex = "1[a-z0-9]+"
width =  450
height = 120
`
