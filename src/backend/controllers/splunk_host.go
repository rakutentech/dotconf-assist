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

func AddSplunkHostHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var splunkHost models.SplunkHost //here is single object, not slice

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &splunkHost); err != nil {
		goto Error
	}

	splunkHost.CreatedAt = time.Now()
	if err = models.SaveSplunkHost(splunkHost); err != nil {
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

func GetSplunkHostsHandler(w http.ResponseWriter, r *http.Request) { //get all splunkHosts
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkHosts []models.SplunkHost

	env := r.URL.Query().Get("env")
	role := r.URL.Query().Get("role")
	token := r.Header.Get("x-auth-token")
	params := []string{ConvertEmpty2Percent(env), ConvertEmpty2Percent(role)}
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if splunkHosts, err = models.GetSplunkHosts(params); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, splunkHosts); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetSplunkHostHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkHost models.SplunkHost

	vars := mux.Vars(r)
	splunkHostName := vars["name"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if splunkHost, err = models.GetSplunkHost(splunkHostName); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, splunkHost); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateSplunkHostHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkHost models.SplunkHost //here is single object, not slice
	var body []byte

	vars := mux.Vars(r)
	splunkHostName := vars["name"]

	token := r.Header.Get("x-auth-token")
	msg, email, _ := ValidateToken(token)
	if email == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(email) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}
	if err = json.Unmarshal(body, &splunkHost); err != nil {
		goto Error
	}

	splunkHost.UpdatedAt = time.Now()
	if err = models.UpdateSplunkHost(splunkHostName, splunkHost); err != nil {
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

func DeleteSplunkHostHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	splunkHostName := vars["name"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if err = models.DeleteSplunkHost(splunkHostName); err != nil {
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
