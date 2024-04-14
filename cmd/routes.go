package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", app.registerUser)
	mux.HandleFunc("POST /login", app.loginUser)

	dynamic := alice.New(app.authenticate)
	mux.Handle("GET /user_banner", dynamic.ThenFunc(app.getBanner))

	protected := dynamic.Append(app.authorization)
	mux.Handle("GET /banner", protected.ThenFunc(app.getAllBanners))
	mux.Handle("POST /banner", protected.ThenFunc(app.createBanner))
	mux.Handle("DELETE /banner/{id}", protected.ThenFunc(app.deleteBanner))
	mux.Handle("PATCH /banner/{id}", protected.ThenFunc(app.patchBanner))

	standard := alice.New(app.recoverPanic, app.logRequest, app.setContentTypeJSON)

	return standard.Then(mux)
}

// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMzNzQzMTcsInJvbGUiOiJhZG1pbiIsInN1YiI6ImJvc3MifQ.AajvkV0kPcb6QuP3wi2jyBxEBLIyl-iDVRVxq-f_VIA
// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMzNzQ0NzcsInJvbGUiOiJ1c2VyIiwic3ViIjoidXNlciJ9.JsEzFOHOThsljtpzr17qe7_XsbUXDyZwURCDUsOzwMU
