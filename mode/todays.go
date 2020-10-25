package mode

import (
	"fmt"
	"strings"
	"time"
)

func ToDays(table string, partKey string) string {
	parts := make([]string, ForwardCreateDays)
	date := time.Now()
	for i := 0; i < ForwardCreateDays; i++ {
		date = date.AddDate(0, 0, 1)
		parts[i] = fmt.Sprintf("PARTITION %s VALUES LESS THAN (TO_DAYS('%s'))", partNameGen(date), date.Format("2006-01-02"))
	}
	return fmt.Sprintf("ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))(%s);", table, partKey, strings.Join(parts, ","))
}
