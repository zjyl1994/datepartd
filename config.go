package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

var conf Config

type Config struct {
	Database []Database
}

type Database struct {
	Name      string
	DSN       string
	Mode      string
	Table     string
	Partkey   string
	PurgeDays int
}

func loadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &conf)
}
