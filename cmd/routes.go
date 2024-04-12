package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.recoverPanic)
	r.Use(app.logRequest)

	r.Get("/user_banner", app.getBanner)
	r.Get("/banner", app.GetAllBanners)
	return r
}
