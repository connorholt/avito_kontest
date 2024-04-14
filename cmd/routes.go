package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	chain := alice.New(app.recoverPanic, app.logRequest, app.setContentTypeJSON)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user_banner", app.getBanner)
	mux.HandleFunc("GET /banner", app.getAllBanners)
	mux.HandleFunc("POST /banner", app.createBanner)
	mux.HandleFunc("DELETE /banner/{id}", app.deleteBanner)
	mux.HandleFunc("PATCH /banner/{id}", app.patchBanner)
	return chain.Then(mux)
}
