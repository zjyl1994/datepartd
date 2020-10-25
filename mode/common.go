package mode

import (
	"database/sql"
	"fmt"
	"time"

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
