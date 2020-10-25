package mode

import (
	"fmt"
	"strings"
)

func ToDaysCreate(mode Database) error {
	parts := make([]string, ForwardCreateDays)
	date := nowTime()
	for i := 0; i < ForwardCreateDays; i++ {
		date = date.AddDate(0, 0, 1)
		parts[i] = fmt.Sprintf("PARTITION %s VALUES LESS THAN (TO_DAYS('%s'))", partNameGen(date), date.Format("2006-01-02"))
	}
	sqlStr := fmt.Sprintf("ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))(%s);", mode.Table, mode.Partkey, strings.Join(parts, ","))
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

func ToDaysDelete(mode Database) error {
	db, err := connectDB(mode.Server, mode.Port, mode.Username, mode.Password, mode.Database)
	if err != nil {
		return err
	}
	sqlTemplate := "SELECT partition_name FROM information_schema.`PARTITIONS` WHERE" +
		` table_schema="%s" AND TABLE_NAME="%s" AND partition_method="RANGE" AND` +
		` PARTITION_EXPRESSION="%s" AND CAST(Partition_description AS UNSIGNED) < TO_DAYS('%s');`
	expression := "to_days(`" + strings.ToLower(mode.Partkey) + "`)"
	purgeDay := nowTime().AddDate(0, 0, -mode.PurgeDays).Format("2006-01-02")
	sqlString := fmt.Sprintf(sqlTemplate, mode.Database, mode.Table, expression, purgeDay)

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
		sqlString = fmt.Sprintf(`ALTER TABLE %s DROP PARTITION %s;`, mode.Table, strings.Join(partNames, ","))
		_, err = db.Exec(sqlString)
		if err != nil {
			return err
		}
	}
	return db.Close()
}
