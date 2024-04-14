package main

import (
	"avito_app/internal/models"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

//	@title			avito_app
//	@version		1.0
//	@description	My Application
//  @schemes 		http
//	@host			localhost:5000
//	@BasePath		/
//  @securityDefinitions.apikey ApiKeyAuth
//  @in header
//  @name Authorization

type application struct {
	infoLog   *log.Logger
	errorLog  *log.Logger
	banners   *models.BannerModel
	bannerTag *models.BannerTagModel
	users     *models.UserModel
	cache     cachedData
}

func main() {
	dbConfig := models.GetConfig()

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
		infoLog:   infoLog,
		errorLog:  errorLog,
		banners:   &models.BannerModel{DB: db},
		bannerTag: &models.BannerTagModel{DB: db},
		users:     &models.UserModel{DB: db},
		cache:     cachedData{},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}
	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
