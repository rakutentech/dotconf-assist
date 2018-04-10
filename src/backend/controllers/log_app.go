package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"strconv"
	"time"
)

func AddAppHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var app models.App //here is single object, not slice
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &app); err != nil {
		goto Error
	}

	if err = models.SaveApp(app); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetAppsHandler(w http.ResponseWriter, r *http.Request) { //get all apps
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var apps []models.App
	env := r.URL.Query().Get("env")
	user := r.URL.Query().Get("user")
	params := []string{ConvertEmpty2Percent(env), ConvertEmpty2Percent(user)}
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	isAdmin := false

	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if userName == "admin" {
		isAdmin = true
	}

	if apps, err = models.GetApps(params, true, true, true, isAdmin); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, apps); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetAppHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var app models.App
	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	params := []string{"", ""} // need to improve
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	isAdmin := false
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if userName == "admin" {
		isAdmin = true
	}

	if app, err = models.GetApp(params, appID, false, isAdmin); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, app); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateAppHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var app models.App //here is single object, not slice
	var body []byte

	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}
	if err = json.Unmarshal(body, &app); err != nil {
		goto Error
	}

	app.UpdatedAt = time.Now()
	if err = models.UpdateApp(appID, app); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func DeleteAppHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if err = models.DeleteApp(appID); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}
