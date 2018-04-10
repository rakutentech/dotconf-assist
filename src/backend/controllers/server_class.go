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

func AddServerClassHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var serverClass models.ServerClass //here is single object, not slice
	var fwdr interface{}
	var f interface{}

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

	if err = json.Unmarshal(body, &f); err != nil {
		goto Error
	}

	serverClass.Name = f.(map[string]interface{})["name"].(string)
	serverClass.UserName = f.(map[string]interface{})["user_name"].(string)
	serverClass.Env = f.(map[string]interface{})["env"].(string)
	for _, fwdr = range f.(map[string]interface{})["forwarders"].([]interface{}) {
		serverClass.ForwarderIDs = append(serverClass.ForwarderIDs, int(fwdr.(map[string]interface{})["id"].(float64)))
	}

	serverClass.CreatedAt = time.Now()
	if err = models.SaveServerClass(serverClass); err != nil {
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

func GetServerClassesHandler(w http.ResponseWriter, r *http.Request) { //get all serverClasses
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var serverClasses []models.ServerClass
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

	if serverClasses, err = models.GetServerClasses(params, isAdmin); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, serverClasses); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetServerClassHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var serverClass models.ServerClass

	vars := mux.Vars(r)
	serverClassName := vars["name"]
	user := r.URL.Query().Get("user")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if serverClass, err = models.GetServerClass(serverClassName, user); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, serverClass); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateServerClassHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var serverClass models.ServerClass //here is single object, not slice
	var fwdr interface{}
	var f interface{}

	vars := mux.Vars(r)
	serverClassName := vars["name"]
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

	if err = json.Unmarshal(body, &f); err != nil {
		goto Error
	}

	serverClass.Name = f.(map[string]interface{})["name"].(string)
	serverClass.UserName = f.(map[string]interface{})["user_name"].(string)
	serverClass.Env = f.(map[string]interface{})["env"].(string)
	for _, fwdr = range f.(map[string]interface{})["forwarders"].([]interface{}) {
		serverClass.ForwarderIDs = append(serverClass.ForwarderIDs, int(fwdr.(map[string]interface{})["id"].(float64)))
	}

	serverClass.UpdatedAt = time.Now()
	if err = models.UpdateServerClass(serverClassName, user, serverClass); err != nil {
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

func DeleteServerClassHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	serverClassName := vars["name"]
	user := r.URL.Query().Get("user")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if err = models.DeleteServerClass(serverClassName, user); err != nil {
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
