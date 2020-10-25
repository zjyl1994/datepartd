package main

import (
	"fmt"
	"log"

	"github.com/zjyl1994/datepartd/mode"
)

var modes = map[string]func(string, string) string{
	"todays": mode.ToDays,
}

func main() {
	err := loadConfig("config.toml")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("%#v\n", conf)

	for _, v := range conf.Database {
		if createFn, ok := modes[v.Mode]; ok {
			sql := createFn(v.Table, v.Partkey)
			fmt.Println(v.Name, "Create SQL", sql)
		}
		if v.PurgeDays > 0 {
			sql := mode.PartDelete(v.Table, v.PurgeDays)
			fmt.Println(v.Name, "Delete SQL", sql)
		}
	}
}
