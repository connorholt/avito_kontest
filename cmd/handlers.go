package main

import (
	"avito_app/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (app *application) getBanner(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	tagId, err := strconv.Atoi(v.Get("tag_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	featureId, err := strconv.Atoi(v.Get("feature_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	banner, err := app.banners.Get(tagId, featureId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data, err := json.MarshalIndent(banner, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	featureID, tagID, limitID, offset := -1, -1, -1, -1
	var err error
	tag := v.Get("tag_id")
	if tag != "" {
		tagID, err = strconv.Atoi(tag)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}
	feature := v.Get("feature_id")
	if feature != "" {
		featureID, err = strconv.Atoi(feature)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}
	limit := v.Get("limit")
	if limit != "" {
		limitID, err = strconv.Atoi(limit)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}
	offsetStr := v.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}
	banners, err := app.banners.GetAll(tagID, featureID, limitID, offset)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data, err := json.MarshalIndent(banners, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
