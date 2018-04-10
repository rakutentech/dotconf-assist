package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	// "strconv"
)

func ResetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var hashedPassword []byte
	var user models.User
	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &user); err != nil {
		goto Error
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		goto Error
	}
	user.Password = string(hashedPassword)

	if err = models.ResetUserPassword(user); err != nil {
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

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var hashedPassword []byte
	var user models.User
	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &user); err != nil {
		goto Error
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		goto Error
	}
	user.Password = string(hashedPassword)
	user.Status = "Waiting"
	if user.UserName == "admin" {
		user.Status = "Approved"
	}
	// i, err := strconv.Atoi("-42")
	if err = models.SaveUser(user, true); err != nil {
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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) { //get all users
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var users []models.User

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	// if !models.IsAdmin(userName) {
	// 	err = errors.New("Permission denied")
	// 	responseCode = http.StatusForbidden
	// 	goto Error
	// }

	if users, err = models.GetUsers(); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, users); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var user models.User

	vars := mux.Vars(r)
	name := vars["username"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) && name != userName {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if user, err = models.GetUser(name); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, user); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var user models.User //here is single object, not slice
	var body []byte
	var hashedPassword []byte
	action := r.URL.Query().Get("action")
	vars := mux.Vars(r)
	name := vars["username"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if !models.IsAdmin(userName) && name != userName {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if !models.IsAdmin(userName) && action == "change_status" {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}
	if err = json.Unmarshal(body, &user); err != nil {
		goto Error
	}
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		goto Error
	}
	user.Password = string(hashedPassword)
	if err = models.UpdateUser(name, user, action); err != nil {
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	name := vars["username"]

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

	if models.IsAdmin(name) {
		err = errors.New("Can't delete Admin")
		responseCode = http.StatusForbidden
		goto Error
	}

	if err = models.DeleteUser(name); err != nil {
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
