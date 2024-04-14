package main

import (
	"avito_app/internal/models"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
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
	var lastRevision bool
	if headValue := r.Header.Get("use_last_revision"); headValue == "true" {
		lastRevision = true
	}
	key := pair{tagID: tagId, featureID: featureId}
	var banner *models.Banner
	value, ok := app.cache[key]
	if ok && !lastRevision && time.Since(value.timestamp) <= 5*time.Minute {
		banner = value.banner
	} else {
		res, err := app.banners.Get(tagId, featureId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, http.StatusNotFound)
			} else {
				app.serverError(w, err)
			}
			return
		}
		banner = res
		app.cache[key] = cacheValue{banner, time.Now()}
	}

	if *banner.IsActive == false {
		app.clientError(w, http.StatusNotFound)
		return
	}
	content := banner.Content
	data, err := json.MarshalIndent(content, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) getAllBanners(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	featureID, tagID, limitID, offset := 0, 0, 0, 0
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
	for _, value := range []int{tagID, featureID, limitID, offset} {
		if value < 0 {
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

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) deleteBanner(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.banners.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createBanner(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = banner.Validate()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if banner.IsActive == nil {
		flag := true
		banner.IsActive = &flag
	}
	bannerID, err := app.banners.Create(banner)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.bannerTag.Create(bannerID, banner.FeatureID, banner.TagID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	res := struct {
		BannerID int `json:"banner_id"`
	}{BannerID: bannerID}
	data, err := json.MarshalIndent(res, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (app *application) patchBanner(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var banner models.Banner
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		app.serverError(w, err)
		return
	}
	banner.ID = id
	err = app.banners.Update(banner)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}
	if len(banner.TagID) != 0 {
		err = app.bannerTag.Delete(banner.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.bannerTag.Create(banner.ID, banner.FeatureID, banner.TagID)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {

	u := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if u.Role == "" {
		u.Role = "user"
	} else if u.Role != "user" && u.Role != "admin" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.users.Insert(u.Name, u.Password, u.Role)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUserName) {
			app.clientError(w, http.StatusBadRequest)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	logReq := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		app.serverError(w, err)
	}

	u, err := app.users.Get(logReq.Name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusUnauthorized)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(logReq.Password)); err != nil {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	payload := jwt.MapClaims{
		"sub":  u.Name,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"role": u.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data, err := json.MarshalIndent(LoginToken{AccessToken: t}, "", "	")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
