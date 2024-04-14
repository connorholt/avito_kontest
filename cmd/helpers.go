package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *application) clientError(w http.ResponseWriter, Status int) {
	http.Error(w, http.StatusText(Status), Status)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func fillTable(table string, db *sql.DB) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("INSERT INTO %s (name) VALUES ", table))
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
