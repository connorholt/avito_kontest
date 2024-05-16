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

// @Summary getBanner
// @ID get-banner-id
// @Param feature_id query integer true "feature_id"
// @Param tag_id query integer true "tag_id"
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /user_banner [get]
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

// @Summary getAllBanners
// @ID get-all-banners-id
// @Param feature_id query integer false "feature_id"
// @Param tag_id query integer false "tag_id"
// @Param limit query integer false "limit"
// @Param offset query integer false "offset"
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal server error"
// @Router /banner [get]
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

// @Summary deleteBanner
// @ID delete-banner-id
// @Param id path integer true "banner id"
// @Security ApiKeyAuth
// @Accept json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /banner/{id} [delete]
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

// @Summary createBanner
// @ID create-banner-id
// @Param input body  models.Banner true "Data about banner"
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal server error"
// @Router /banner [post]
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

// @Summary patchBanner
// @ID patch-banner-id
// @Param id path integer true "banner id"
// @Param input body  models.Banner true "Data about banner"
// @Security ApiKeyAuth
// @Accept json
// @Success 200 "OK"
// @Failure 400 "Client Error"
// @Failure 401 "You are not authorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal server error"
// @Router /banner/{id} [patch]
func (app *application) updateBanner(w http.ResponseWriter, r *http.Request) {
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

// @Summary registerUser
// @Description register new user
// @Tags registration and logging
// @Param input body RegisterRequest true "New user info"
// @Accept json
// @Produce json
// @success 201 {integer} integer "New user registered"
// @Failure 409 {object} error
// @Router /register [post]
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

// @Summary loginUser
// @Description login existing user
// @Tags registration and logging
// @Param input body LoginRequest true "New user info"
// @Accept json
// @Produce json
// @success 200 {integer} integer "Succesfuly logged in
// @Failure 409 {object} error
// @Router /login [post]
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

// @Summary fillTables
// @Description fill tables features and tags 1000 rows
// @Tags preparing
// @success 200 {integer} integer
// @Failure 500 "Internal server error"
// @Router /fill_tables [post]
func (app *application) fillFeatureAndTags(w http.ResponseWriter, r *http.Request) {
	err := fillTable("features", app.banners.DB)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = fillTable("tags", app.banners.DB)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
