package main

import (
	"flag"
	"fmt"

	"avito_app/internal/models"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	banners  *models.BannerModel
}

func main() {
	dbConfig := GetConfig()

	addr := flag.String("addr", ":5000", "HTTP network address")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DB_name)

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	db, err := models.OpenDB(psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		banners:  &models.BannerModel{db},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}
	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
