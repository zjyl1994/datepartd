package mode

import (
	"fmt"
	"time"
)

const (
	ForwardCreateDays = 7
)

func partNameGen(date time.Time) string {
	return "p" + date.Format("20060102")
}

func PartDelete(table string, purgeDays int) string {
	partName := partNameGen(time.Now().AddDate(0, 0, -purgeDays))
	return fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s;", table, partName)
}
