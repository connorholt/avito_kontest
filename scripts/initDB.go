package scripts

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func FillTable(table string, db *sql.DB) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("INSERT INTO %s VALUES ", table))
	for i := 0; i < 1000; i++ {
		b.WriteString(fmt.Sprintf("('%s'), ", strconv.Itoa(i)))
	}
	b.WriteString("('1000')")
	stmt := b.String()
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
