package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("DATEPARTD_START")
	err := loadConfig("config.toml")
	if err != nil {
		log.Println(err.Error())
		return
	}
	createPartitionJob()
	deletePartitionJob()
	// c := cron.New()
	// _, err = c.AddFunc("0 23 * * *", createPartitionJob)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// _, err = c.AddFunc("0 1 * * *", deletePartitionJob)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// c.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	<-chSignal
	log.Println("DATEPARTD_STOP")
}

func createPartitionJob() {
	for _, v := range conf.Database {
		err := createFn(v)
		log.Println(v.Name, "Create", err)
	}
}

func deletePartitionJob() {
	for _, v := range conf.Database {
		if v.PurgeDays > 0 {
			err := deleteFn(v)
			log.Println(v.Name, "Delete", err)
		}
	}
}
