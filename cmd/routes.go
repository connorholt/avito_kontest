package main

import (
	_ "avito_app/docs"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", app.registerUser)
	mux.HandleFunc("POST /login", app.loginUser)
	mux.HandleFunc("POST /fill_tables", app.fillFeatureAndTags)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	dynamic := alice.New(app.authenticate)
	mux.Handle("GET /user_banner", dynamic.ThenFunc(app.getBanner))

	protected := dynamic.Append(app.authorization)
	mux.Handle("GET /banner", protected.ThenFunc(app.getAllBanners))
	mux.Handle("POST /banner", protected.ThenFunc(app.createBanner))
	mux.Handle("DELETE /banner/{id}", protected.ThenFunc(app.deleteBanner))
	mux.Handle("PATCH /banner/{id}", protected.ThenFunc(app.updateBanner))

	standard := alice.New(app.recoverPanic, app.logRequest, app.setContentTypeJSON)

	return standard.Then(mux)
}
