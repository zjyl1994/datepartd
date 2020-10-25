package main

import (
	"log"

	"github.com/robfig/cron/v3"
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
	c := cron.New()
	_, err = c.AddFunc("0 23 * * *", createPartitionJob)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = c.AddFunc("0 1 * * *", deletePartitionJob)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.Start()

}

func createPartitionJob() {
	for _, v := range conf.Database {
		if createFn, ok := createModes[v.Mode]; ok {
			err := createFn(v)
			log.Println(v.Name, "Create", err)
		}
	}
}

func deletePartitionJob() {
	for _, v := range conf.Database {
		if v.PurgeDays > 0 {
			if deleteFn, ok := deleteModes[v.Mode]; ok {
				err := deleteFn(v)
				log.Println(v.Name, "Delete", err)
			}
		}
	}
}
