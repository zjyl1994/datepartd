package main

import (
	"log"
	"os"
	"os/signal"

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
	log.Println("DATEPARTD_START")
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
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	<-chSignal
	log.Println("DATEPARTD_STOP")
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
