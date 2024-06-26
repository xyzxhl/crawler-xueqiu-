package db

import (
	"CrawlBase/pj"
	"bytes"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MultInit struct{}

func (MultInit) Error() string {
	return "Multiple initializations"
}

var (
	db *sql.DB
	mi MultInit
)

func InitDB() error {
	if db != nil {
		return mi
	}

	tmp, err := sql.Open("mysql", "root:12345678@tcp(47.120.8.50:3306)/Funds_Indices?parseTime=true")
	if err != nil {
		return err
	}

	db = tmp
	return nil
}

func FINameInsert(symbol string, name string) error {
	_, err := db.Exec("INSERT INTO FIName(symbol,name,earliest_date)VALUES(?,?,?)", symbol, name, nil)
	return err
}

func FINameUpdateDate(dates map[string]time.Time) error {
	instBase := "INSERT INTO FIName(symbol,earliest_date)VALUES"
	var buf bytes.Buffer
	buf.WriteString(instBase)
	vals := []interface{}{}
	for symbol, date := range dates {
		buf.WriteString("(?,?),")
		vals = append(vals, symbol, date)
	}
	inst := buf.String()
	inst = inst[:len(inst)-1]
	inst = inst + " ON DUPLICATE KEY UPDATE earliest_date=VALUES(earliest_date)"

	_, err := db.Exec(inst, vals...)
	return err
}

func FINameGetAllSymbols() ([]string, error) {
	symbols := make([]string, 0)
	rows, err := db.Query("SELECT symbol FROM FIName")
	if err != nil {
		return nil, err
	}

	var tmp string
	for rows.Next() {
		if err := rows.Scan(&tmp); err != nil {
			return nil, err
		}
		symbols = append(symbols, tmp)
	}

	return symbols, nil
}

func CHRecordInsert(data map[string][]pj.LChange) error {
	instBase := "INSERT INTO CHRecord(symbol,date,percent)VALUES"
	var buf bytes.Buffer
	buf.WriteString(instBase)
	vals := []interface{}{}
	for symbol, changes := range data {
		for _, v := range changes {
			buf.WriteString("(?,?,?),")
			vals = append(vals, symbol, v.T, v.Percent)
		}
	}
	inst := buf.String()
	inst = inst[:len(inst)-1]

	_, err := db.Exec(inst, vals...)
	return err
}
