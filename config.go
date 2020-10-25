package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/zjyl1994/datepartd/mode"
)

var conf Config

type Config struct {
	Database []mode.Database
}

func loadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &conf)
}
