package main

import (
	"database/sql"
	"fmt"
	"time"

	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Name      string
	Server    string
	Port      int
	Username  string
	Password  string
	Database  string
	Mode      string
	Table     string
	Partkey   string
	PurgeDays int
}

const (
	ForwardCreateDays = 7
)

func partNameGen(date time.Time) string {
	return "p" + date.Format("20060102")
}

func connectDB(server string, port int, username, password, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, server, port, database)
	return sql.Open("mysql", dsn)
}

func nowTime() time.Time {
	return time.Now()
}

func createFn(mode Database) error {
	parts := make([]string, ForwardCreateDays)
	date := nowTime()
	for i := 0; i < ForwardCreateDays; i++ {
		date = date.AddDate(0, 0, 1)
		parts[i] = fmt.Sprintf("PARTITION %s VALUES LESS THAN (%s('%s'))", partNameGen(date), strings.ToUpper(mode.Mode), date.Format("2006-01-02"))
	}
	sqlStr := fmt.Sprintf("ALTER TABLE %s PARTITION BY RANGE (%s(%s))(%s);", mode.Table, strings.ToUpper(mode.Mode), mode.Partkey, strings.Join(parts, ","))
	log.Println("CREATE_PARTITION_SQL", sqlStr)
	db, err := connectDB(mode.Server, mode.Port, mode.Username, mode.Password, mode.Database)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlStr)
	if err != nil {
		return err
	}
	return db.Close()
}

func deleteFn(mode Database) error {
	db, err := connectDB(mode.Server, mode.Port, mode.Username, mode.Password, mode.Database)
	if err != nil {
		return err
	}
	sqlTemplate := "SELECT partition_name FROM information_schema.`PARTITIONS` WHERE" +
		` table_schema="%s" AND TABLE_NAME="%s" AND partition_method="RANGE" AND` +
		` LOWER(PARTITION_EXPRESSION)="%s" AND CAST(Partition_description AS UNSIGNED) < %s('%s');`
	expression := strings.ToLower(mode.Mode) + "(`" + strings.ToLower(mode.Partkey) + "`)"
	purgeDay := nowTime().AddDate(0, 0, -mode.PurgeDays).Format("2006-01-02")
	sqlString := fmt.Sprintf(sqlTemplate, mode.Database, mode.Table, expression, strings.ToUpper(mode.Mode), purgeDay)
	log.Println("QUERY_PARTITION_SQL", sqlString)
	rows, err := db.Query(sqlString)
	if err != nil {
		return err
	}
	partNames := make([]string, 0)
	for rows.Next() {
		var partName string
		err = rows.Scan(&partName)
		if err != nil {
			return err
		}
		partNames = append(partNames, partName)
	}
	err = rows.Close()
	if err != nil {
		return err
	}
	if len(partNames) > 0 {
		log.Println("PARTITION_WILL_DELETE", strings.Join(partNames, ","))
		sqlString = fmt.Sprintf(`ALTER TABLE %s DROP PARTITION %s;`, mode.Table, strings.Join(partNames, ","))
		log.Println("DELETE_PARTITION_SQL", sqlString)
		_, err = db.Exec(sqlString)
		if err != nil {
			return err
		}
	} else {
		log.Println("NO_PARTITION_WILL_DELETE")
	}
	return db.Close()
}
