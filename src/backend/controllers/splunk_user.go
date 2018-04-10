package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	// "golang.org/x/crypto/bcrypt"
	"net/http"
	// "strconv"
)

func AddSplunkUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var splunkSplunkUser models.SplunkUser //here is single object, not slice

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &splunkSplunkUser); err != nil {
		goto Error
	}

	ConcertNull2EmptyStr(&splunkSplunkUser.Memo)
	ConcertNull2EmptyStr(&splunkSplunkUser.RpaasUserName)

	if err = models.SaveSplunkUser(splunkSplunkUser); err != nil {
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

func GetSplunkUsersHandler(w http.ResponseWriter, r *http.Request) { //get all splunkSplunkUsers
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkSplunkUsers []models.SplunkUser

	token := r.Header.Get("x-auth-token")
	msg, splunkSplunkUserName, _ := ValidateToken(token)
	if splunkSplunkUserName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(splunkSplunkUserName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if splunkSplunkUsers, err = models.GetSplunkUsers(); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, splunkSplunkUsers); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetSplunkUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkSplunkUser models.SplunkUser

	env := r.URL.Query().Get("env")
	vars := mux.Vars(r)
	name := vars["username"]

	token := r.Header.Get("x-auth-token")
	msg, splunkSplunkUserName, _ := ValidateToken(token)
	if splunkSplunkUserName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(splunkSplunkUserName) && name != splunkSplunkUserName {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if splunkSplunkUser, err = models.GetSplunkUser(name, ConvertEmpty2Percent(env)); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, splunkSplunkUser); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateSplunkUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var splunkSplunkUser models.SplunkUser //here is single object, not slice
	var body []byte

	vars := mux.Vars(r)
	name := vars["username"]

	token := r.Header.Get("x-auth-token")
	msg, splunkSplunkUserName, _ := ValidateToken(token)
	if splunkSplunkUserName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(splunkSplunkUserName) && name != splunkSplunkUserName {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}
	if err = json.Unmarshal(body, &splunkSplunkUser); err != nil {
		goto Error
	}

	if err = models.UpdateSplunkUser(name, splunkSplunkUser); err != nil {
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

func DeleteSplunkUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	name := vars["username"]

	token := r.Header.Get("x-auth-token")
	msg, splunkSplunkUserName, _ := ValidateToken(token)
	if splunkSplunkUserName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(splunkSplunkUserName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if models.IsAdmin(name) {
		err = errors.New("Can't delete Admin")
		responseCode = http.StatusForbidden
		goto Error
	}

	if err = models.DeleteSplunkUser(name); err != nil {
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
