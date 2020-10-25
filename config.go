package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

var conf Config

type Config struct {
	Cron     bool
	DryRun   bool
	Timezone string
	Database []Database
}

func loadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &conf)
}
