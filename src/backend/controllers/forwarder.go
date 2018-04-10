package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"time"
)

func AddForwarderHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var forwarders []models.Forwarder

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

	if err = json.Unmarshal(body, &forwarders); err != nil {
		goto Error
	}

	// forwarder.CreatedAt = time.Now()
	if err = models.SaveForwarder(forwarders); err != nil {
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

func GetForwardersHandler(w http.ResponseWriter, r *http.Request) { //get all forwarders
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var forwarders []models.Forwarder
	env := r.URL.Query().Get("env")
	user := r.URL.Query().Get("user")
	from := r.URL.Query().Get("from")
	params := []string{ConvertEmpty2Percent(env), ConvertEmpty2Percent(user), ConvertEmpty2Percent(from)}
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

	if forwarders, err = models.GetForwarders(params, isAdmin); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, forwarders); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetForwarderHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var forwarder models.Forwarder

	vars := mux.Vars(r)
	forwarderName := vars["name"]
	user := r.URL.Query().Get("user")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if forwarder, err = models.GetForwarder(forwarderName, user); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, forwarder); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateForwarderHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var forwarder models.Forwarder //here is single object, not slice
	var body []byte

	vars := mux.Vars(r)
	forwarderName := vars["name"]
	user := r.URL.Query().Get("user")

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
	if err = json.Unmarshal(body, &forwarder); err != nil {
		goto Error
	}

	forwarder.UpdatedAt = time.Now()
	if err = models.UpdateForwarder(forwarderName, user, forwarder); err != nil {
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

func DeleteForwarderHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	forwarderName := vars["name"]
	user := r.URL.Query().Get("user")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if err = models.DeleteForwarder(forwarderName, user); err != nil {
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
