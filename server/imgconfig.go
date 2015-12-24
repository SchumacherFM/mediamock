package server

import (
	"regexp"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Dirs []dir
}

type dir struct {
	regex  *regexp.Regexp
	Name   string
	Regex  string
	Width  int
	Height int
}

func parseImgConfigFile(path string) (*tomlConfig, error) {

	var cfg = new(tomlConfig)
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, err
	}

	for _, c := range cfg.Dirs {
		if c.Regex != "" {
			var err error
			if c.regex, err = regexp.Compile(c.Regex); err != nil {

			}
		}
	}

	return cfg, nil
}
