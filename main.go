package main

import (
	"fmt"
	"log"

	"github.com/zjyl1994/datepartd/mode"
)

var createModes = map[string]func(mode.Database) error{
	"todays": mode.ToDaysCreate,
}

var deleteModes = map[string]func(mode.Database) error{
	"todays": mode.ToDaysDelete,
}

func main() {
	err := loadConfig("config.toml")
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, v := range conf.Database {
		if createFn, ok := createModes[v.Mode]; ok {
			err := createFn(v)
			fmt.Println(v.Name, "Create", err)
		}
		if v.PurgeDays > 0 {
			if deleteFn, ok := deleteModes[v.Mode]; ok {
				err := deleteFn(v)
				fmt.Println(v.Name, "Delete", err)
			}
		}
	}
}
